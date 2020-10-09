[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/arl/statsviz)

Statsviz
========

Instant live visualization of your Go application runtime statistics 
(GC, MemStats, etc.).

 - Import `import _ "github.com/arl/statsviz"` (à la `"net/http/pprof"`)
 - Open your browser at `http://host:port/debug/statsviz`
 - Enjoy... 


Installation
------------

```bash
go get -u github.com/arl/statsviz
```


Usage
-----

This package is typically only imported for the side effect of registering its
HTTP handler. The handled path is `/debug/statsviz/`.

If your application is not already running an HTTP server, you need to start
one. Add `"net/http"` and `"log"` to your imports and the following code to your
`main` function:

```go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

If you are not using [http.DefaultServeMux](https://pkg.go.dev/net/http?tab=doc#ServeMux),
you will have to register the handler with the mux you are using.

Then open your browser at http://localhost:6060/debug/statsviz/


Plots
-----

On the plots where it matters, garbage collections are shown as vertical bars.
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
More details in [Contributing](CONTRIBUTING.md)


License
-------

- [MIT License](LICENSE)
