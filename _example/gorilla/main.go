package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create a Gorilla router.
	r := mux.NewRouter()

	// Create statsviz endpoint and register the handlers on the router.
	ep := statsviz.NewEndpoint()
	r.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(ep.Ws())
	r.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(ep.Index())

	mux := http.NewServeMux()
	mux.Handle("/", r)

	fmt.Println("Point your browser to http://localhost:8086/debug/statsviz/")
	http.ListenAndServe(":8086", mux)
}
