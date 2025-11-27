package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/redis/go-redis/v9"
)

// User struct for data mapping
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB        // Global database connection pool
var rdb *redis.Client // Global Redis client

func initRedis() {
	var ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis!")
}

func initDB() {
	// IMPORTANT: Replace with your actual credentials
	connStr := "user=postgres password=yourpassword dbname=mydb host=localhost sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection is alive
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

//

func main() {
	initDB()
	defer db.Close()

	initRedis()
	defer rdb.Close()

	e := echo.New()

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	// Simple API with database access
	e.GET("/users/:id", getUser)
	// Endpoint for simulating the deadlock
	e.POST("/transfer_deadlock", simulateDeadlockHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// Simple API handler example
func getUser(c echo.Context) error {
	id := c.Param("id")
	var user User

	// 1. Check Redis cache first
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s:name", id)
	cachedName, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		userID, _ := strconv.Atoi(id)
		user.ID = userID
		user.Name = cachedName
		return c.JSON(http.StatusOK, user)
	}

	// 2. Cache miss, fetch from database
	err = db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// 3. Store in Redis cache for future requests with a TTL of 10 minutes
	err = rdb.Set(ctx, cacheKey, user.Name, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set cache for user %s: %v", id, err)
	}

	return c.JSON(http.StatusOK, user)
}

// simulateDeadlockHandler handles the API call to start the simulation.
func simulateDeadlockHandler(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Transaction 1: Lock account 1, then try to lock account 2
	go func() {
		defer wg.Done()
		err := runTransaction(1, 2, 50, "Transaction 1")
		if err != nil {
			log.Printf("Tx 1 FAILED: %v", err)
		}
	}()

	// Give Tx 1 a moment to acquire its first lock
	time.Sleep(100 * time.Millisecond)

	// Transaction 2: Lock account 2, then try to lock account 1
	go func() {
		defer wg.Done()
		err := runTransaction(2, 1, 50, "Transaction 2")
		if err != nil {
			log.Printf("Tx 2 FAILED: %v", err)
		}
	}()

	wg.Wait()
	return c.String(http.StatusOK, "Deadlock simulation started. Check logs for result.")
}

// runTransaction attempts to transfer amount from 'fromID' to 'toID'.
func runTransaction(fromID, toID, amount int, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to start transaction: %w", name, err)
	}
	defer tx.Rollback() // Rollback is safe to call even if Commit succeeds

	log.Printf("%s: Started. Acquiring lock on account %d.", name, fromID)

	// STEP 1: Lock 'from' account (Resource 1)
	// Use SELECT FOR UPDATE to acquire a row-level lock
	_, err = tx.Exec(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, fromID)
	if err != nil {
		return fmt.Errorf("%s: failed to lock account %d: %w", name, fromID, err)
	}

	// Introduce a delay here to ensure the other transaction can start and acquire its first lock,
	// creating the circular dependency (the heart of the deadlock).
	time.Sleep(500 * time.Millisecond)

	log.Printf("%s: Acquired lock on %d. Attempting to acquire lock on account %d.", name, fromID, toID)

	// STEP 2: Lock 'to' account (Resource 2)
	// This is where the deadlock occurs: Tx 1 waits for Tx 2's lock on account 2,
	// while Tx 2 (in its first step) is waiting for Tx 1's lock on account 1.
	_, err = tx.Exec(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, toID)
	if err != nil {
		// PostgreSQL detects the deadlock and one transaction will fail here with a '40P01' error.
		return fmt.Errorf("%s: **DEADLOCK POINT** failed to lock account %d: %w", name, toID, err)
	}

	// Perform the actual update/transfer (simplified)
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return fmt.Errorf("%s: failed to update balance for %d: %w", name, fromID, err)
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return fmt.Errorf("%s: failed to update balance for %d: %w", name, toID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", name, err)
	}

	log.Printf("%s: SUCCESSFULLY committed transfer of %d from %d to %d.", name, amount, fromID, toID)
	return nil
}
