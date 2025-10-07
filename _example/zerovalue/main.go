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

	// Register with the Server zero value.
	var ss statsviz.Server
	ss.Register(http.DefaultServeMux)

	fmt.Println("Point your browser to http://localhost:8079/debug/statsviz/")
	log.Fatal(http.ListenAndServe(":8079", nil))
}
