// Package statsviz serves via its HTTP server an HTML page displaying live
// visualization of the application runtime statistics.
//
// The package is typically only imported for the side effect of
// registering its HTTP handler.
// The handled path is /debug/statsviz/.
//
// To use statsviz, link this package into your program:
//	import _ "github.com/arl/statsviz"
//
// If your application is not already running an http server, you
// need to start one. Add "net/http" and "log" to your imports and
// the following code to your main function:
//
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
//
// If you are not using DefaultServeMux, you will have to register handlers
// with the mux you are using.
//
// Then open your browser at http://localhost:6060/debug/statsviz/
package statsviz

import (
	"net/http"
)

// Register registers statsviz HTTP handlers on the provided mux.
func Register(mux *http.ServeMux) {
	mux.Handle("/debug/statsviz/", Index)
	mux.HandleFunc("/debug/statsviz/ws", Ws)
}

// RegisterDefault registers statsviz HTTP handlers on the default serve mux.
//
// Note this is not advised on a production server, unless it only serves on
// localhost.
func RegisterDefault() {
	Register(http.DefaultServeMux)
}
