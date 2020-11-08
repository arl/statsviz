package main

import (
	"log"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create a serve mux and register statsviz handlers at /foo/bar
	mux := http.NewServeMux()

	if err := statsviz.Register(mux, statsviz.Root("/foo/bar")); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", mux))
}
