package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

// basicAuth adds HTTP Basic Authentication to h.
//
// NOTE: This is just an example middleware to show how one can wrap statsviz
// handler, it should absolutely not be used as is.
func basicAuth(h http.HandlerFunc, user, pwd, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); !ok || user != u || pwd != p {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		h(w, r)
	}
}

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	const statsvizRoot = "/debug/statsviz"

	// Create statsviz endpoint.
	se := statsviz.NewEndpoint()

	mux := http.NewServeMux()
	mux.Handle(statsvizRoot+"/", basicAuth(se.Index(), "statsviz", "rocks", ""))
	mux.HandleFunc(statsvizRoot+"/ws", se.Ws())

	fmt.Println("Point your browser to http://localhost:8090/debug/statsviz/")
	fmt.Println("Basic auth user:     statsviz")
	fmt.Println("Basic auth password: rocks")
	log.Fatal(http.ListenAndServe(":8090", mux))
}
