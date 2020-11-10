package main

import (
	"github.com/fasthttp/websocket"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Fiber instance
	app := fiber.New()

	// Use fiber adaptor to wrap statsviz handler
	app.Get("/debug/statsviz", adaptor.HTTPHandler(statsviz.Index))
	// Create fasthttp websocket handler
	app.Get("/debug/statsviz/ws", func(c *fiber.Ctx) error {
		var upgrader = websocket.FastHTTPUpgrader{}
		err := upgrader.Upgrade(c.Context(), func(ws *websocket.Conn) {
			_ = statsviz.SendStats(ws)
		})
		return err
	})

	// Start server
	app.Listen(":8080")
}
