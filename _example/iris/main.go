package main

import (
	"fmt"
	"net/http"

	"github.com/kataras/iris/v12"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	app := iris.New()

	// Need to run iris in a separate goroutine so we can start the dedicated
	// http server for Statsviz.
	go app.Listen(":8089")

	mux := http.NewServeMux()

	// Create and register statsviz endpoint.
	se := statsviz.NewEndpoint()
	se.Register(mux)
	statsSrv := &http.Server{Addr: ":8088", Handler: mux}

	fmt.Println("Point your browser to http://localhost:8088/debug/statsviz\n")

	// NewHost puts the http server for statsviz under the control of iris but
	// iris won't touch its handlers.
	app.NewHost(statsSrv).ListenAndServe()
}
