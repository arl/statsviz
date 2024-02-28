package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
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
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Create statsviz server.
	srv, err := statsviz.NewServer()

	// Register Statsviz server on the fasthttp router.
	app.Use("/debug/statsviz", adaptor.HTTPHandler(srv.Index()))
	ws.HandleFunc("/debug/statsviz/ws", srv.Ws())

	fmt.Println("Point your browser to http://localhost:8093/debug/statsviz/")

	// Server start
	go http.Serve(m.Match(cmux.HTTP1HeaderField("Accept", "text/event-stream")), ws)
	go app.Listener(m.Match(cmux.Any()))
	m.Serve()
}
