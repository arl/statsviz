package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create a serve mux.
	mux := http.NewServeMux()

	// Create statsviz endpoint and registers it with the serve mux.
	se := statsviz.NewEndpoint()
	se.Register(mux)

	fmt.Println("Point your browser to http://localhost:8091/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8091", mux))
}
