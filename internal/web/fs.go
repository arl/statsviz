package web

import (
	"embed"
	"fmt"
	"io/fs"
	"sync"
)

//go:embed dist
var dist embed.FS

var Dist = sync.OnceValue(func() fs.FS {
	var err error
	webFS, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(fmt.Sprintf("fs failed to get 'dist' subtree: %s", err))
	}
	return webFS
})
