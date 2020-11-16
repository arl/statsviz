package main

import (
	"net"
	"net/http"

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
	l, _ := net.Listen("tcp", ":8080")
	m := cmux.New(l)
	ws := http.NewServeMux()

	// Fiber instance
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// statsviz
	app.Get("/debug/statsviz", adaptor.HTTPHandler(statsviz.Index))
	ws.HandleFunc("/debug/statsviz/ws", statsviz.Ws)

	// Server start
	go http.Serve(m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket")), ws)
	go app.Listener(m.Match(cmux.Any()))
	m.Serve()
}
