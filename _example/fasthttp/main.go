package main

import (
	"fmt"
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

	// Create the main listener and mux
	l, _ := net.Listen("tcp", ":8083")
	m := cmux.New(l)
	ws := http.NewServeMux()

	// fasthttp routers
	r := router.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world!")
	})

	// Create statsviz endpoint.
	se := statsviz.NewEndpoint()

	r.GET("/debug/statsviz/{filepath:*}", fasthttpadaptor.NewFastHTTPHandler(se.Index()))
	ws.HandleFunc("/debug/statsviz/ws", se.Ws())

	// Server start
	go http.Serve(m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket")), ws)
	go fasthttp.Serve(m.Match(cmux.Any()), r.Handler)
	fmt.Println("Point your browser to http://localhost:8083/debug/statsviz/")
	m.Serve()
}
