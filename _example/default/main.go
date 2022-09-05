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

	// Register statsviz handlers on the default serve mux.
	statsviz.RegisterDefault()

	fmt.Println("Point your browser to http://localhost:8080/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
