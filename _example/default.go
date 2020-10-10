package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/arl/statsviz"
)

func garbage() []byte {
	var b []byte

	rnd := rand.New(rand.NewSource(0))
	switch rnd.Intn(4) {
	case 0:
		b = make([]byte, 8192+rnd.Intn(8192*4))
	case 1:
		b = make([]byte, 2048+rnd.Intn(4096))
	case 2, 3:
		b = make([]byte, rnd.Intn(128))
	}

	return b
}

func main() {
	// Force the GC to work
	go func() {
		m := make(map[string][]byte)
		i := 0
		for {
			m[fmt.Sprintf("%d", i)] = garbage()
			time.Sleep(10 * time.Millisecond)
			i++
			if i%(10*100) == 0 {
				m = make(map[string][]byte)
			}
		}
	}()

	// Register statsviz handlers on the default serve mux.
	statsviz.RegisterDefault()
	http.ListenAndServe(":8080", nil)
}
