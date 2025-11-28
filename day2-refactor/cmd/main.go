package main

import (
	"context"
	"demo"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var conn *pgxpool.Pool // Global database connection pool

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

type StubRepository struct{}

func (r *StubRepository) GetByID(ctx context.Context, id string) (*demo.User, error) {
	fmt.Println("Called stub ...")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	if userID == 5 {
		return nil, pgx.ErrNoRows
	}
	return &demo.User{ID: userID, Name: "Stub User"}, nil
}

func main() {

	user := demo.UserHandler{
		Repo: &StubRepository{},
	}

	e := echo.New()

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	// Simple API with database access
	e.GET("/users/:id", user.GetUser)

	e.Logger.Fatal(e.Start(":8080"))
}
