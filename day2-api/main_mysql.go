package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var db *sql.DB // Global database connection pool
var stmt *sql.Stmt

func initDB() {
	var err error
	db, err = sql.Open("mysql", "user:yourpassword@/mydb")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	// Test the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Prepare the statement once
	stmt, err = db.Prepare("SELECT content FROM messages WHERE id = ?")
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to mysql !!")

}

//

func main() {
	initDB()
	defer db.Close()
	defer stmt.Close()

	e := echo.New()

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Get message by id from mysql
	e.GET("/messages/:id", func(c echo.Context) error {
		id := c.Param("id")
		var message string
		err := stmt.QueryRowContext(c.Request().Context(), id).Scan(&message)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.String(http.StatusNotFound, "Message not found")
			}
			return c.String(http.StatusInternalServerError, "Database error")
		}
		return c.String(http.StatusOK, message)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
