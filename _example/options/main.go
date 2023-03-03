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

	// Serve statsviz user interface from /foo/bar instead of /debug/statsviz
	// and send metrics every 250 milliseconds instead of 1 second.
	ep := statsviz.NewEndpoint(
		statsviz.WithRoot("/foo/bar"),
		statsviz.WithInterval(250*time.Millisecond),
	)

	mux := http.NewServeMux()

	// Register the endpoint handlers on the mux.
	ep.Register(mux)

	log.Println("Point your browser to http://localhost:8092/foo/bar")
	log.Fatal(http.ListenAndServe(":8092", mux))
}
