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

	mux := http.NewServeMux()

	// Register the endpoint handlers on the mux.
	statsviz.NewEndpoint().Register(mux)

	fmt.Println("Point your browser to http://localhost:8091/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8091", mux))
}
