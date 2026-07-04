package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/soheilhy/cmux"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create the main listener and mux
	l, err := net.Listen("tcp", ":8093")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m := cmux.New(l)
	ws := http.NewServeMux()

	// Fiber instance
	app := fiber.New()

	// Create statsviz server.
	srv, err := statsviz.NewServer()
	if err != nil {
		panic(err)
	}

	app.Use("/debug/statsviz/", srv.Index())
	ws.HandleFunc("/debug/statsviz/ws", srv.Ws())

	fmt.Println("Point your browser to http://localhost:8093/debug/statsviz/")

	// Server start
	go http.Serve(m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket")), ws)
	go app.Listener(m.Match(cmux.Any()))
	m.Serve()
}
