package main

import (
	"log"
	"net/http"
	"time"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create statsviz endpoint at path /foo/bar with a metrics send interval of 250ms.
	se := statsviz.NewEndpoint().
		WithRoot("/foo/bar").
		WithSendInterval(250 * time.Millisecond)

	// Create a serve mux and register statsviz endpoint.
	mux := http.NewServeMux()
	se.Register(mux)

	log.Println("Point your browser to http://localhost:8092/foo/bar")
	log.Fatal(http.ListenAndServe(":8092", mux))
}
