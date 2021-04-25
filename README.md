[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/arl/statsviz)
[![Test Actions Status](https://github.com/arl/statsviz/workflows/Test/badge.svg)](https://github.com/arl/statsviz/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/arl/statsviz)](https://goreportcard.com/report/github.com/arl/statsviz)
[![codecov](https://codecov.io/gh/arl/statsviz/branch/master/graph/badge.svg)](https://codecov.io/gh/arl/statsviz)


Statsviz
========

Instant live visualization of your Go application runtime statistics 
(GC, MemStats, etc.).

 - Import `"github.com/arl/statsviz"`
 - Register statsviz HTTP handlers
 - Start your program
 - Open your browser at `http://host:port/debug/statsviz`
 - Enjoy... 


How does that work?
-----------------

Statsviz serves 2 HTTP handlers.

The first one (by default `/debug/statsviz`) serves an html/js user interface showing 
some initially empty plots.

When you points your browser to statsviz user interface page, it connects to statsviz
second HTTP handler. This second handler then upgrades the connection to the websocket
protocol and starts a goroutine that periodically calls [runtime.ReadMemStats](https://golang.org/pkg/runtime/#ReadMemStats), 
sending the result to the user interface, which inturn, updates the plots.

Stats are stored in-browser inside a circular buffer which keep tracks of a predefined number of
datapoints, 60, so one minute-worth of data, by default. You can change the frequency
at which stats are sent by passing [SendFrequency](https://pkg.go.dev/github.com/arl/statsviz@v0.2.1#SendFrequency)
to [Register](https://pkg.go.dev/github.com/arl/statsviz@v0.2.1#Register).


Usage
-----

    go get github.com/arl/statsviz

Either `Register` statsviz HTTP handlers with the [http.ServeMux](https://pkg.go.dev/net/http?tab=doc#ServeMux) you're using (preferred method):

```go
mux := http.NewServeMux()
statsviz.Register(mux)
```

Or register them with the `http.DefaultServeMux`:

```go
statsviz.RegisterDefault()
```

If your application is not already running an HTTP server, you need to start
one. Add `"net/http"` and `"log"` to your imports and the following code to your
`main` function:

```go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

By default the handled path is `/debug/statsviz/`.

Then open your browser at http://localhost:6060/debug/statsviz/

Plots
-----

On the plots where it matters, garbage collections are shown as vertical lines.

### Heap
<img alt="Heap plot image" src="https://github.com/arl/statsviz/raw/readme-docs/heap.png" width="600">

### MSpans / MCaches
<img alt="MSpan/MCache plot image" src="https://github.com/arl/statsviz/raw/readme-docs/mspan-mcache.png" width="600">

### Size classes heatmap
<img alt="Size classes heatmap image" src="https://github.com/arl/statsviz/raw/readme-docs/size-classes.png" width="600">

### Objects
<img alt="Objects plot image" src="https://github.com/arl/statsviz/raw/readme-docs/objects.png" width="600">

### Goroutines
<img alt="Goroutines plot image" src="https://github.com/arl/statsviz/raw/readme-docs/goroutines.png" width="600">

### GC/CPU fraction
<img alt="GC/CPU fraction plot image" src="https://github.com/arl/statsviz/raw/readme-docs/gc-cpu-fraction.png" width="600">


Examples
--------

Have a look at the [_example](./_example/README.md) directory to see some
different ways to register Statsviz HTTP handlers, such as:
 - using `http.DefaultServeMux`
 - using your own `http.ServeMux`
 - register at `/foo/bar` instead of `/debug/statviz`
 - use `https://` rather than `http://`
 - using with various Go HTTP libraries/frameworks:
   - [fasthttp](https://github.com/valyala/fasthttp)
   - [gin](https://github.com/gin-gonic/gin)
   - and many others thanks to wonderful contributors!


Contributing
------------

Pull-requests are welcome!
More details in [CONTRIBUTING.md](CONTRIBUTING.md)


Roadmap
-------

 - [ ] add stop-the-world duration heatmap
 - [ ] increase data retention
 - [ ] light/dark mode selector
 - [x] plot image export as png
 - [ ] save timeseries to disk
 - [ ] load from disk previously saved timeseries


Changelog
---------

See [CHANGELOG.md](./CHANGELOG.md).


License
-------

- [MIT License](LICENSE)