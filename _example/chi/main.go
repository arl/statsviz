package main

import (
	"fmt"
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

	fmt.Println("Point your browser to http://localhost:8081/debug/statsviz/")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
