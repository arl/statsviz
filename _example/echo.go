package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/arl/statsviz"
)

func main(){
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	// Statsviz
	// Add monitor Go application runtime statistics (GC, MemStats, etc.)
	// Use http DefaultServeMux
	mux := http.DefaultServeMux
	// Register route
	_ = statsviz.Register(mux)

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/debug/statsviz/", echo.WrapHandler(mux))
	// Server static content for statsviz UI
	e.GET("/debug/statsviz/*", echo.WrapHandler(mux))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}