// Package statsviz serves via its HTTP server an HTML page displaying live
// visualization of the application runtime statistics.
//
// Either Register statsviz HTTP handlers with the http.ServeMux you're using
// (preferred method):
//  mux := http.NewServeMux()
//  statsviz.Register(mux)
//
// Or register them with the http.DefaultServeMux:
//  statsviz.RegisterDefault()
//
// If your application is not already running an HTTP server, you need to start
// one. Add "net/http" and "log" to your imports and the following code to your
// main function:
//  go func() {
//  	log.Println(http.ListenAndServe("localhost:6060", nil))
//  }()
package statsviz
