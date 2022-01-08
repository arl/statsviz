package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

type ServeStatus struct {
	AvailableChannel int
	ClientCount      int
	MaxChannelCount  int
}
type NetworkStatus struct {
	LinkCount    int   `json:"link count"`
	DataSent     int64 `json:"data sent"`
	DataReceived int64 `json:"data received"`
	Other        int64
}
type CustomData struct {
	NetworkStatus NetworkStatus `json:"Network status"`
	ServeStatus   ServeStatus   `json:"Server status"`
}

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	data := &CustomData{}
	data.ServeStatus.MaxChannelCount = 1000
	//Generate custom data
	statsviz.CustomDataGenerate = func() interface{} {
		data.NetworkStatus.LinkCount = rand.Intn(100)
		data.NetworkStatus.DataSent += rand.Int63n(100)
		data.NetworkStatus.DataReceived += rand.Int63n(100)
		data.ServeStatus.ClientCount = rand.Intn(1000)
		data.ServeStatus.AvailableChannel = rand.Intn(500)
		return data
	}

	// Register statsviz handlers on the default serve mux.
	statsviz.RegisterDefault()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
