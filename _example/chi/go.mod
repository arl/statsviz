module example/chi

go 1.23

toolchain go1.24.5

require (
	github.com/arl/statsviz v0.7.0
	github.com/go-chi/chi v1.5.4
)

require github.com/gorilla/websocket v1.5.3 // indirect

replace github.com/arl/statsviz => ../../
