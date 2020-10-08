// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// This program uses `github.com/shurcooL/vfsgen` to generate the `assets`
// variable. `assets` statically implements an `http.FileSystem` rooted at
// `/static/`, it embeds the files that statsviz serves.
//
// While working on statsviz web interface, it's easier to directly serve the
// content of the /static directory rather than regenerating the assets after
// each modification. Passing `-tags dev` to `go build` will do just that, the
// directory served will be that of your filesystem.
//
// However to commit the modifications of the /static directory, `assets` must
// be regenerated, to do so just call `go generate` from the project root.
// With Go modules enabled, this will download the latest version of
// github.com/shurcooL/vfsgen and update `assets_vfsdata.go` so that it
// reflects the new content of the /static directory. Then, commits both
// /static and assets_vfsdata.go.

func main() {
	err := vfsgen.Generate(http.Dir("static"), vfsgen.Options{
		PackageName:  "statsviz",
		BuildTags:    "!dev",
		VariableName: "assets",
		Filename:     "assets_vfsdata.go",
	})
	if err != nil {
		log.Fatalln("generate assets", err)
	}
}
