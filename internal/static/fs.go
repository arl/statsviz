package static

import (
	"embed"
	"fmt"
	"io/fs"
	"sync"
)

//go:embed dist
var assets embed.FS

var Assets = sync.OnceValue(func() fs.FS {
	var err error
	webFS, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(fmt.Sprintf("error loading frontend assets: %s", err))
	}
	return webFS
})
