// Package statsviz allows to visualise Go program runtime metrics data in real
// time: heap, objects, goroutines, GC pauses, scheduler, etc. in your browser.
//
// Register statsviz endpoint on your server http.ServeMux (preferred method):
//
//	mux := http.NewServeMux()
//	statsviz.Register(mux)
//
// Or register on `http.DefaultServeMux`:
//
//	statsviz.RegisterDefault()
//
// By default statsviz is served at `/debug/statsviz/`.
//
// If your application is not already running an HTTP server, you need to start
// one. Add "net/http" and "log" to your imports and the following code to your
// main function:
//
//	go func() {
//	    log.Println(http.ListenAndServe("localhost:6060", nil))
//	}()
//
// Then open your browser at http://localhost:6060/debug/statsviz/.
package statsviz
