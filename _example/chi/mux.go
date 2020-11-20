package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/arl/statsviz"
	"github.com/go-chi/chi"
)

func main() {

	// Force the GC to work to make the plots "move".
	go work()

	// Create a Chi router and register statsviz handlers.
	r := chi.NewRouter()
	r.Get("/debug/statsviz/ws", statsviz.Ws)
	r.Handle("/debug/statsviz/*", statsviz.Index)

	mux := http.NewServeMux()
	mux.Handle("/", r)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}

func work() {
	// Generate some allocations
	m := map[string][]byte{}

	for {
		b := make([]byte, 512+rand.Intn(16*1024))
		m[strconv.Itoa(len(m)%(10*100))] = b

		if len(m)%(10*100) == 0 {
			m = make(map[string][]byte)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
