module example/chi

go 1.19

require (
	github.com/arl/statsviz v0.6.0
	github.com/go-chi/chi v1.5.4
)

require github.com/gorilla/websocket v1.5.0 // indirect

replace github.com/arl/statsviz => ../../
