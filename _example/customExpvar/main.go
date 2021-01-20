package main

import (
	"expvar"
	"log"
	"math/rand"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func init() {
	expvar.Publish("customInt", expvar.Func(func() interface{} {
		return rand.Intn(100)
	}))
	expvar.Publish("customFloat", expvar.Func(func() interface{} {
		return rand.Float64() * 100.0
	}))
	e := expvar.NewMap("customMap")
	e.Add("a", 1)
	e.Add("b", 5)
}

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	// Register statsviz handlers on the default serve mux.
	statsviz.RegisterDefault()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
