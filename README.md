[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/arl/statsviz)
[![Test Actions Status](https://github.com/arl/statsviz/workflows/Test/badge.svg)](https://github.com/arl/statsviz/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/arl/statsviz)](https://goreportcard.com/report/github.com/arl/statsviz)
[![codecov](https://codecov.io/gh/arl/statsviz/branch/master/graph/badge.svg)](https://codecov.io/gh/arl/statsviz)

Statsviz
========

Instant live visualization of your Go application runtime statistics 
(GC, MemStats, etc.).

 - Import `import "github.com/arl/statsviz"`
 - Register statsviz HTTP handlers
 - Start your program
 - Open your browser at `http://host:port/debug/statsviz`
 - Enjoy... 


Usage
-----

    go get -u github.com/arl/statsviz

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

Examples
--------

Using `http.DefaultServeMux`:
 - [_example/default.go](./_example/default.go)

Using your own `http.ServeMux`:
 - [_example/mux.go](./_example/mux.go)

Using https`:
 - [_example/https.go](./_example/https.go)

With [gorilla/mux](https://github.com/gorilla/mux) Router:
 - [_example/gorilla/mux.go](./_example/gorilla/mux.go)


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


Contributing
------------

Pull-requests are welcome!
More details in [CONTRIBUTING.md](CONTRIBUTING.md)


Roadmap
-------

 - [ ] add stop-the-world duration heatmap
 - [ ] increase data retention
 - [ ] light/dark mode selector
 - [ ] plot image export

Changelog
---------

See [CHANGELOG.md](./CHANGELOG.md).

License
-------

- [MIT License](LICENSE)
