package main

import (
	"net"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/soheilhy/cmux"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create the main listener and cmux
	l, _ := net.Listen("tcp", ":8080")
	m := cmux.New(l)

	// Fasthttp router
	r := router.New()
	r.GET("/fasthttp/example", func(ctx *fasthttp.RequestCtx) {})
	// statsviz http
	r.GET("/debug/statsviz/{filepath:*}", fasthttpadaptor.NewFastHTTPHandler(statsviz.Index))

	// statsviz websocket
	ws := http.NewServeMux()
	ws.HandleFunc("/debug/statsviz/ws", statsviz.Ws)

	// Server start
	go http.Serve(m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket")), ws)
	go fasthttp.Serve(m.Match(cmux.Any()), r.Handler)
	m.Serve()
}
