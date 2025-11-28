package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	msg := "Hello, World!"

	e := echo.New()
	e.GET("/", hello(msg))
	e.Logger.Fatal(e.Start(":8080"))
}

func hello(message string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, message)
	}
}
