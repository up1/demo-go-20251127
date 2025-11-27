package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/redis/go-redis/v9"
)

// User struct for data mapping
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var conn *pgxpool.Pool // Global database connection pool
var rdb *redis.Client  // Global Redis client

func initRedis() {
	var ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Working with pool settings
	rdb.Options().PoolSize = 25                     // Maximum number of connections
	rdb.Options().MinIdleConns = 10                 // Minimum number of idle connections
	rdb.Options().MaxIdleConns = 25                 // Maximum number of idle connections
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
	var err error
	connStr := "postgres://postgres:yourpassword@localhost:5432/mydb?sslmode=disable"
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}
	poolConfig.MaxConns = 20                     // Maximum number of connections
	poolConfig.MinConns = 5                      // Minimum number of connections
	poolConfig.MaxConnIdleTime = 5 * time.Minute // Maximum idle time

	conn, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	// Test connection
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

//

func main() {
	initDB()
	defer conn.Close()

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

	// Automatically prepared and reused statement
	err := conn.QueryRow(context.Background(), "SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
