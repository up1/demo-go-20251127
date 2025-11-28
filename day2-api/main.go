package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type resource struct {
	message string
}

func main() {
	r := resource{message: "Hello, World v2!"}

	e := echo.New()
	e.GET("/", hello("V1 Hello, World!"))
	e.GET("/v2", r.hello2)
	e.Logger.Fatal(e.Start(":8080"))
}

// Solution 2
func (r resource) hello2(c echo.Context) error {
	return c.String(http.StatusOK, r.message)
}

// Solution 1
func hello(message string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, message)
	}
}
