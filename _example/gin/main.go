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

	fmt.Println("Point your browser to http://localhost:8085/debug/statsviz/\n\n")

	router := gin.New()
	router.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			statsviz.Ws(context.Writer, context.Request)
			return
		}
		statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(context.Writer, context.Request)
	})
	router.Run(":8085")
}
