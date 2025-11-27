package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	// Working with pool settings
	rdb.Options().PoolSize = 20                     // Maximum number of connections
	rdb.Options().MinIdleConns = 10                 // Minimum number of idle connections
	rdb.Options().ConnMaxIdleTime = 5 * time.Minute // Maximum idle time

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
