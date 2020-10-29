package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create a Gorilla router and register statsviz handlers.
	r := mux.NewRouter()
	r.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(statsviz.Ws)
	r.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(statsviz.Index)

	mux := http.NewServeMux()
	mux.Handle("/", r)
	http.ListenAndServe(":8080", mux)
}
