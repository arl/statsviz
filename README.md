[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=round-square)](https://pkg.go.dev/github.com/arl/statsviz)
[![Latest tag](https://img.shields.io/github/tag/arl/statsviz.svg)](https://github.com/arl/statsviz/tag/)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[![Test Actions Status](https://github.com/arl/statsviz/workflows/Tests-linux/badge.svg)](https://github.com/arl/statsviz/actions)
[![Test Actions Status](https://github.com/arl/statsviz/workflows/Tests-others/badge.svg)](https://github.com/arl/statsviz/actions)
[![codecov](https://codecov.io/gh/arl/statsviz/branch/main/graph/badge.svg)](https://codecov.io/gh/arl/statsviz)

# Statsviz

<p align="center">
  <img alt="Statsviz Gopher Logo" width="120" src="https://raw.githubusercontent.com/arl/statsviz/readme-docs/logo.png?sanitize=true">
  <img alt="statsviz ui" width="450" align="right" src="https://github.com/arl/statsviz/raw/readme-docs/window.png">
</p>
<br/>

Visualize real time plots of your Go program runtime metrics, including heap, objects, goroutines, GC pauses, scheduler and more, in your browser.

<hr>

- [Statsviz](#statsviz)
  - [Install](#install)
  - [Usage](#usage)
  - [Advanced Usage](#advanced-usage)
  - [How Does That Work?](#how-does-that-work)
  - [Documentation](#documentation)
    - [Go API](#go-api)
    - [Web User Interface](#web-user-interface)
    - [Plots](#plots)
    - [User Plots](#user-plots)
  - [Examples](#examples)
  - [Questions / Troubleshooting](#questions--troubleshooting)
  - [Contributing](#contributing)
  - [Changelog](#changelog)
  - [License: MIT](#license-mit)

## Install

Download the latest version:

```
go get github.com/arl/statsviz@latest
```

Please note that, as new metrics are added to the `/runtime/metrics` package, new plots are added to Statsviz.
This also means that the presence of some plots on the dashboard depends on the Go version you're using.

When in doubt, use the latest ;-)


## Usage

Register `Statsviz` HTTP handlers with your application `http.ServeMux`.

```go
mux := http.NewServeMux()
statsviz.Register(mux)

go func() {
    log.Println(http.ListenAndServe("localhost:8080", mux))
}()
```

Open your browser at http://localhost:8080/debug/statsviz


## Advanced Usage

If you want more control over Statsviz HTTP handlers, examples are:
 - you're using some HTTP framework
 - you want to place Statsviz handler behind some middleware

then use `statsviz.NewServer` to obtain a `Server` instance. Both the `Index()` and `Ws()` methods return `http.HandlerFunc`.

```go
srv, err := statsviz.NewServer(); // Create server or handle error
srv.Index()                       // UI (dashboard) http.HandlerFunc
srv.Ws()                          // Websocket http.HandlerFunc
```

Please look at examples of usage in the [Examples](_example) directory.


## How Does That Work?

`statsviz.Register` registers 2 HTTP handlers within the given `http.ServeMux`:

- the `Index` handler serves Statsviz user interface at `/debug/statsviz` at the address served by your program.

- The `Ws` serves a Websocket endpoint. When the browser connects to that endpoint, [runtime/metrics](https://pkg.go.dev/runtime/metrics) are sent to the browser, once per second.

Data points are in a browser-side circular-buffer.


## Documentation

### Go API

Check out the API reference on [pkg.go.dev](https://pkg.go.dev/github.com/arl/statsviz#section-documentation).

### Web User Interface


#### Top Bar

<img alt="webui-annotated" src="https://github.com/arl/statsviz/raw/readme-docs/webui-annotated.png">

##### Category Selector

<img alt="menu-categories" src="https://github.com/arl/statsviz/raw/readme-docs/menu-categories.png">

Each plot belongs to one or more categories. The category selector allows you to filter the visible plots by categories.

##### Visible Time Range

<img alt="menu-timerange" src="https://github.com/arl/statsviz/raw/readme-docs/menu-timerange.png">

Use the time range selector to define the visualized time span.

##### Show/Hide GC events

<img alt="menu-gc-events" src="https://github.com/arl/statsviz/raw/readme-docs/menu-gc-events.png">

Show or hide the vertical lines representing garbage collection events.

##### Pause updates

<img alt="menu-play" src="https://github.com/arl/statsviz/raw/readme-docs/menu-play.png">

Pause or resume the plot updates.


#### Plot Controls

<img alt="webui-annotated" src="https://github.com/arl/statsviz/raw/readme-docs/plot-controls-annotated.png">


### Plots

The visible set of plots depend on your Go version since some plots are only available in newer versions.

#### Allocation and Free Rate

<img width="50%" alt="alloc-free-rate" src="https://github.com/arl/statsviz/raw/readme-docs/plots/alloc-free-rate.png">

#### CGO Calls

<img width="50%" alt="cgo" src="https://github.com/arl/statsviz/raw/readme-docs/plots/cgo.png">

#### CPU (GC)

<img width="50%" alt="cpu-gc" src="https://github.com/arl/statsviz/raw/readme-docs/plots/cpu-gc.png">

#### CPU (Overall)

<img width="50%" alt="cpu-overall" src="https://github.com/arl/statsviz/raw/readme-docs/plots/cpu-overall.png">

#### CPU (Scavenger)

<img width="50%" alt="cpu-scavenger" src="https://github.com/arl/statsviz/raw/readme-docs/plots/cpu-scavenger.png">

#### Garbage Collection

<img width="50%" alt="garbage-collection" src="https://github.com/arl/statsviz/raw/readme-docs/plots/garbage collection.png">

#### GC Cycles

<img width="50%" alt="gc-cycles" src="https://github.com/arl/statsviz/raw/readme-docs/plots/gc-cycles.png">

#### GC Pauses

<img width="50%" alt="gc-pauses" src="https://github.com/arl/statsviz/raw/readme-docs/plots/gc-pauses.png">

#### GC Scan

<img width="50%" alt="gc-scan" src="https://github.com/arl/statsviz/raw/readme-docs/plots/gc-scan.png">

#### GC Stack Size

<img width="50%" alt="gc-stack-size" src="https://github.com/arl/statsviz/raw/readme-docs/plots/gc-stack-size.png">

#### Goroutines

<img width="50%" alt="goroutines" src="https://github.com/arl/statsviz/raw/readme-docs/plots/goroutines.png">

#### Heap (Details)

<img width="50%" alt="heap-details" src="https://github.com/arl/statsviz/raw/readme-docs/plots/heap (details).png">

#### Live Bytes

<img width="50%" alt="live-bytes" src="https://github.com/arl/statsviz/raw/readme-docs/plots/live-bytes.png">

#### Live Objects

<img width="50%" alt="live-objects" src="https://github.com/arl/statsviz/raw/readme-docs/plots/live-objects.png">

#### Memory Classes

<img width="50%" alt="memory-classes" src="https://github.com/arl/statsviz/raw/readme-docs/plots/memory-classes.png">

#### MSpan/MCache

<img width="50%" alt="mspan-mcache" src="https://github.com/arl/statsviz/raw/readme-docs/plots/mspan-mcache.png">

#### Mutex Wait

<img width="50%" alt="mutex-wait" src="https://github.com/arl/statsviz/raw/readme-docs/plots/mutex-wait.png">

#### Runnable Time

<img width="50%" alt="runnable-time" src="https://github.com/arl/statsviz/raw/readme-docs/plots/runnable-time.png">

#### Scheduling Events

<img width="50%" alt="sched-events" src="https://github.com/arl/statsviz/raw/readme-docs/plots/sched-events.png">

#### Size Classes

<img width="50%" alt="size-classes" src="https://github.com/arl/statsviz/raw/readme-docs/plots/size-classes.png">

#### GC Pauses

<img width="50%" alt="gc-pauses" src="https://github.com/arl/statsviz/raw/readme-docs/plots/gc-pauses.png">


### User Plots

Since `v0.6` you can add your own plots to Statsviz dashboard, in order to easily
visualize your application metrics next to runtime metrics.

Please see the [userplots example](_example/userplots/main.go).

## Examples

Check out the [\_example](./_example/README.md) directory to see various ways to use Statsviz, such as:

- use of `http.DefaultServeMux` or your own `http.ServeMux`
- wrap HTTP handler behind a middleware
- register the web page at `/foo/bar` instead of `/debug/statsviz`
- use `https://` rather than `http://`
- register Statsviz handlers with various Go HTTP libraries/frameworks:
  - [echo](https://github.com/labstack/echo/)
  - [fasthttp](https://github.com/valyala/fasthttp)
  - [fiber](https://github.com/gofiber/fiber/)
  - [gin](https://github.com/gin-gonic/gin)
  - and many others thanks to many contributors!

## Questions / Troubleshooting

Either use GitHub's [discussions](https://github.com/arl/statsviz/discussions) or come to say hi and ask a live question on [#statsviz channel on Gopher's slack](https://gophers.slack.com/archives/C043DU4NZ9D).

## Contributing

Please use [issues](https://github.com/arl/statsviz/issues/new/choose) for bugs and feature requests.  
Pull-requests are always welcome!  
More details in [CONTRIBUTING.md](CONTRIBUTING.md).

## Changelog

See [CHANGELOG.md](./CHANGELOG.md).

## License: MIT

See [LICENSE](LICENSE)
