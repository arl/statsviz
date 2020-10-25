package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/arl/statsviz"
)

func main() {
	// Force the GC to work to make the plots "move".
	go work()

	// Create a serve mux and register statsviz handlers at /foo/bar
	mux := http.NewServeMux()

	if err := statsviz.Register(mux, statsviz.Root("/debug")); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", mux))
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
