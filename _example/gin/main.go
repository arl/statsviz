package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/arl/statsviz"
	example "github.com/arl/statsviz/_example"
)

func main() {
	// Force the GC to work to make the plots "move".
	go example.Work()

	fmt.Printf("Point your browser to http://localhost:8085/debug/statsviz/\n\n")

	// Create statsviz endpoint.
	se := statsviz.NewEndpoint()

	ws := se.Ws()
	index := se.Index()

	router := gin.New()
	router.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			ws(context.Writer, context.Request)
			return
		}
		index(context.Writer, context.Request)
	})
	router.Run(":8085")
}
