Examples
========

## Using [net/http](https://pkg.go.dev/net/http)

Using `http.DefaultServeMux`:
 - [default/main.go](./default/main.go)

Using your own `http.ServeMux`:
 - [mux/main.go](./mux/main.go)

Use statsviz options API to serve Statsviz web UI on `/foo/bar` (instead of default
`/debug/statsviz`) and send metrics with a frequency of _250ms_ (rather than _1s_):
 - [options/main.go](./options/main.go)

Serve the the web UI via `https` and Websocket via `wss`:
 - [https/main.go](./https/main.go)

Wrap statviz handlers behind a middleware (HTTP Basic Authentication for example):
 - [middleware/main.go](./middleware/main.go)


## Using various Go libraries

With [gorilla/mux](https://github.com/gorilla/mux) router:
 - [gorilla/main.go](./gorilla/main.go)

Using [valyala/fasthttp](https://github.com/valyala/fasthttp) and [soheilhy/cmux](https://github.com/soheilhy/cmux):
 - [fasthttp/main.go](./fasthttp/main.go)

Using [labstack/echo](https://github.com/labstack/echo) router:
 - [echo/main.go](./echo/main.go)

With [gin-gonic/gin](https://github.com/gin-gonic/gin) web framework:
 - [gin/main.go](./gin/main.go)

With [go-chi/chi](https://github.com/go-chi/chi) router:
 - [chi/main.go](./chi/main.go)

With [gofiber/fiber](https://github.com/gofiber/fiber) web framework:
 - [fiber/main.go](./fiber/main.go)

With [kataras/iris](https://github.com/kataras/iris) web framework:
 - [iris/main.go](./iris/main.go)

