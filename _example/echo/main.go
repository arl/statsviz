package main

import (
	example "github.com/arl/statsviz/_example"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Echo instance
	e := echo.New()

	mux := http.NewServeMux()
	_ = statsviz.Register(mux)

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/debug/statsviz/", echo.WrapHandler(mux))
	// Serve static content for statsviz UI
	e.GET("/debug/statsviz/*", echo.WrapHandler(mux))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
