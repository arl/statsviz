package main

import (
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

	mux := http.NewServeMux()
	mux.Handle(statsvizRoot+"/", basicAuth(statsviz.IndexAtRoot(statsvizRoot), "hello", "world", ""))
	mux.HandleFunc(statsvizRoot+"/ws", statsviz.Ws)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
