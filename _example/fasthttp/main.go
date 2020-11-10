package main

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Fasthttp router
	r := router.New()

	// Use fasthttpadaptor to wrap statsviz handler
	r.GET("/debug/statsviz/{filepath:*}", fasthttpadaptor.NewFastHTTPHandler(statsviz.Index))
	// Create fasthttp websocket handler
	r.GET("/debug/statsviz/ws", func(ctx *fasthttp.RequestCtx) {
		var upgrader = websocket.FastHTTPUpgrader{}
		err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
			_ = statsviz.SendStats(ws)
		})
		if err != nil {
			log.Println(err)
			return
		}
	})

	// Start server
	fasthttp.ListenAndServe(":8080", r.Handler)
}
