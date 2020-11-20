package main

import (
	"log"
	"net/http"

	"github.com/arl/statsviz"
	"github.com/go-chi/chi"

	example "github.com/arl/statsviz/_example"
)

func main() {

	// Force the GC to work to make the plots "move".
	go example.Work()

	// Create a Chi router and register statsviz handlers.
	r := chi.NewRouter()
	r.Get("/debug/statsviz/ws", statsviz.Ws)
	r.Get("/debug/statsviz", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/debug/statsviz/", 301)
	})
	r.Handle("/debug/statsviz/*", statsviz.Index)

	mux := http.NewServeMux()
	mux.Handle("/", r)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
