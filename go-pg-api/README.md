# Workshop REST API with Go
* [Echo framework](https://echo.labstack.com/)
* [PostgreSQL database](https://github.com/lib/pq)
* Caching with [Redis](https://github.com/redis/go-redis)

## 1. Create project
* create a new folder = `api`
```
$cd api
$go mod init api
```

## 2. REST API Spexification

API Specification
```
GET /users/:id

Response code = 200
{
    "id": 1,
    "name": "Name 01"
}

Response code = 404
{
    "message": "User not found"
}
```

## 3. Create simple program in file `main.go`
```
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// User struct for data mapping
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB // Global database connection pool

func initDB() {
	// IMPORTANT: Replace with your actual credentials
	connStr := "user=postgres password=yourpassword dbname=deadlock_db host=localhost sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection is alive
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Create a simple table for the deadlock simulation
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)`)
	if err != nil {
		log.Fatal(err)
	}
	// Initial data (if tables are empty)
	db.Exec(`INSERT INTO accounts (id, balance) VALUES (1, 100), (2, 200) ON CONFLICT (id) DO NOTHING`)
}

//

func main() {
	initDB()
	defer db.Close()

	e := echo.New()

	// Simple API endpoint
	e.GET("/users/:id", getUser)

	e.Logger.Fatal(e.Start(":8080"))
}

// Simple API handler example
func getUser(c echo.Context) error {
	id := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "User not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
```

## 4. Run in development mode
```
# Start database
$docker compose up -d
$docker compose ps

# Run server
$go mod tidy
$go run main.go
```

List of URLs
* GET http://localhost:8080/users/1
* GET http://localhost:8080/users/3

