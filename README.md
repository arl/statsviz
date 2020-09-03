# Easymon

Instant live visualization of your Go application runtime statistics (GC,
MemStats, etc.).

 - Only depends on Go standard library.
 - Import `import _ "github.com/arl/easymon"` (à la `"net/http/pprof"`)
 - Open your browser at `http://host:port/debug/easymon`
 - Enjoy...!


## Usage

This package is typically only imported for the side effect of registering its
HTTP handler. The handled path is `/debug/easymon/`.

If your application is not already running an HTTP server, you need to start
one. Add `"net/http"` and `"log"` to your imports and the following code to your
`main` function:

```go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

If you are not using http.DefaultServeMux, you will have to register the handler
with the mux you are using.
