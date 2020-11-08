package main

import (
	"net/http"

	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()

	// Statsviz
	// Add monitor Go application runtime statistics (GC, MemStats, etc.)
	// Create a new http ServeMux
	mux := http.NewServeMux()
	// Register route
	_ = statsviz.Register(mux)

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/debug/statsviz/", echo.WrapHandler(mux))
	// Server static content for statsviz UI
	e.GET("/debug/statsviz/*", echo.WrapHandler(mux))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
