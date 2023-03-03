package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Echo instance
	e := echo.New()

	mux := http.NewServeMux()

	// Create statsviz endpoint and register it with the mux.
	statsviz.NewEndpoint().Register(mux)

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/debug/statsviz/", echo.WrapHandler(mux))
	// Serve static content for statsviz UI
	e.GET("/debug/statsviz/*", echo.WrapHandler(mux))

	// Start server
	fmt.Println("Point your browser to http://localhost:8082/debug/statsviz/")
	e.Logger.Fatal(e.Start(":8082"))
}
