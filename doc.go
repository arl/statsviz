// Package statsviz allows to visualise Go program runtime metrics data in real
// time: heap, objects, goroutines, GC pauses, scheduler, etc. in your browser.
//
// Create a statsviz [Endpoint] and register it with your server [http.ServeMux]
// (preferred method):
//
//	mux := http.NewServeMux()
//	endpoint := statvis.NewEndpoint()
//	endpoint.Register(mux)
//
// Or register with [http.DefaultServeMux`]:
//
//	endpoint := statvis.NewEndpoint()
//	endpoint.Register(http.DefaultServeMux)
//
// By default Statsviz is served at `/debug/statsviz/`. You can change that (and
// other things) using methods on the [statsviz.Endpoint] instance.
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
