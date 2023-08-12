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

	mux := http.NewServeMux()

	// Register Statsviz server on the mux, serving the user interface from
	// /foo/bar instead of /debug/statsviz and send metrics every 250
	// milliseconds instead of the default of once per second.
	_ = statsviz.Register(mux,
		statsviz.Root("/foo/bar"),
		statsviz.SendFrequency(250*time.Millisecond),
	)

	log.Println("Point your browser to http://localhost:8092/foo/bar")
	log.Fatal(http.ListenAndServe(":8092", mux))
}
