package statsviz

import "net/http"

// Index responds to a request for /debug/statsviz with the statsviz HTML page
// which shows a live visualization of the statistics sent by the application
// over the websocket handler Ws.
//
// The package initialization registers it as /debug/statsviz/.
var Index = http.StripPrefix("/debug/statsviz/", http.FileServer(assets))
