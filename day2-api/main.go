package main

import (
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	// Swagger
	_ "api/docs" // Import generated swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

type resource struct {
	message string
}

// Builder pattern => creational
func NewResource(message string) resource {
	return resource{message: message}
}

// Title: Swagger API Example
// @version 1.0
// @description This is a sample server for a Swagger API example.
// @host localhost:8080
// @BasePath /
func main() {
	r := NewResource("Hello, World v2!")
	e := echo.New()
	// e.Use(middleware.Logger())

	// Add prometheus middleware
	e.Use(echoprometheus.NewMiddleware("myapp"))   // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	e.GET("/", hello("V1 Hello, World!"))
	e.GET("/v2", r.hello2)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// Hello2 API
// @Summary Hello2 endpoint
// @Description Returns a hello message from resource
// @Tags hello
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Hello, World v2!"
// @Router /v2 [get]
func (r resource) hello2(c echo.Context) error {
	return c.String(http.StatusOK, r.message)
}

// Hello API
// @Summary Hello endpoint
// @Description Returns a hello message
// @Tags hello
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func hello(message string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, message)
	}
}
