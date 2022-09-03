# Statsviz

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=round-square)](https://pkg.go.dev/github.com/arl/statsviz)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Latest tag](https://img.shields.io/github/tag/arl/statsviz.svg)](https://github.com/arl/statsviz/tag/)  


[![Test Actions Status](https://github.com/arl/statsviz/workflows/Tests-linux/badge.svg)](https://github.com/arl/statsviz/actions)
[![Test Actions Status](https://github.com/arl/statsviz/workflows/Tests-others/badge.svg)](https://github.com/arl/statsviz/actions)
[![codecov](https://codecov.io/gh/arl/statsviz/branch/main/graph/badge.svg)](https://codecov.io/gh/arl/statsviz)

<p align="center">
  <img alt="Statsviz Gopher Logo" width="160" src="https://raw.githubusercontent.com/arl/statsviz/readme-docs/logo.png?sanitize=true">
</p>
<br />

Visualise Go program runtime metrics data in real time: heap, objects, goroutines, GC pauses, scheduler, etc.

 - Import `"github.com/arl/statsviz"`
 - Register statsviz HTTP handlers
 - Start your program
 - Open your browser at `http://host:port/debug/statsviz`
 - Enjoy... 


## Usage

Download the latest version of the `statsviz` module.

    go get github.com/arl/statsviz@latest


Register statsviz HTTP handlers on [http.ServeMux](https://pkg.go.dev/net/http?tab=doc#ServeMux) (preferred method):

```go
mux := http.NewServeMux()
statsviz.Register(mux)
```

Or register them with on `http.DefaultServeMux`:

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

Then open your browser at http://localhost:6060/debug/statsviz/.


## How does that work?

Statsviz serves 2 HTTP endpoints:

 - The first one (`/debug/statsviz`) serves a web page with statsviz
user interface, showing initially empty plots.

 - The second HTTP handler (`/debug/statsviz/ws`) listens for a WebSocket
connection that will be initiated by statsviz web page as soon as it's loaded in
your browser.

That's it, now your application sends all [runtime/metrics](https://pkg.go.dev/runtime/metrics) 
data points to the web page, once per second.

Data points are stored in-browser in a circular buffer which keep tracks of a
predefined number of datapoints.


## Documentation

Check out the API reference on [pkg.go.dev](https://pkg.go.dev/github.com/arl/statsviz#section-documentation).


Plots
-----

On the plots where it matters, garbage collections are shown as vertical lines.

### Heap (global)
<img alt="Heap (global) image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/heap-global.png">

### Heap (details)
<img alt="Heap (details) image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/heap-details.png">

### Live Objects in Heap
<img alt="Live Objects in Heap image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/live%20objects.png">

### Live Bytes in Heap
<img alt="Live Bytes in Heap image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/live%20bytes.png">

### MSpan/MCache
<img alt="MSpan/MCache image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/mspan-mcache.png">

### Goroutines
<img alt="Goroutines image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/goroutines.png">

### Size Classes
<img alt="Size Classes image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/sizeclasses.png">

### Stop-the-world Pause Latencies
<img alt="Stop-the-world Pause Latencies image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/gcpauses.png">

### Time Goroutines Spend in 'Runnable'
<img alt="Time Goroutines Spend in 'Runnable' image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/sched-latencies.png">

### Starting Size of Goroutines Stacks
<img alt="Time Goroutines Spend in 'Runnable' image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/gc_stack-size.png">

### Goroutine Scheduling Events
<img alt="Time Goroutines Spend in 'Runnable' image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/sched%20events.png">

### CGO Calls
<img alt="CGO Calls image" src="https://github.com/arl/statsviz/raw/readme-docs/runtime-metrics/cgo.png">


## Examples directory

Have a look at the [_example](./_example/README.md) directory to see various ways to use Statsviz, such as:
 - using `http.DefaultServeMux`
 - using your own `http.ServeMux`
 - wrap HTTP handler behind a middleware
 - register at `/foo/bar` instead of `/debug/statviz`
 - use `https://` rather than `http://`
 - using with various Go HTTP libraries/frameworks:
   - [fasthttp](https://github.com/valyala/fasthttp)
   - [gin](https://github.com/gin-gonic/gin)
   - and many others thanks to awesome contributors!


## Contributing

Pull-requests are welcome!
More details in [CONTRIBUTING.md](CONTRIBUTING.md).


## Changelog

See [CHANGELOG.md](./CHANGELOG.md).


## License

 See [MIT License](LICENSE)
