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

	// Create a serve mux and register statsviz handlers at /foo/bar with a send
	// frequency of 250ms
	mux := http.NewServeMux()

	err := statsviz.Register(mux,
		statsviz.Root("/foo/bar"),
		statsviz.SendFrequency(250*time.Millisecond),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Point your browser to http://localhost:8080/foo/bar")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
