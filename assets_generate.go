// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// This proram uses `github.com/shurcooL/vfsgen` to generate the `assets`
// variable. `assets` statically implements an `http.FileSystem` rooted at
// /static/`, and thus contains the files rtprof serves.
//
// Just use the `-dev` build tag when developing to directly use the assets
// in the `./static` directory.
//
// However when commiting a modified asset `./assets_vfsdata.go` must be
// re-generated. To do so:
//
// Ensure you have the latest version of `vfsgen`:
//
//    go get -u github.com/shurcooL/vfsgen
//    go generate							# from the project root
//    git add assets_vfsdata.go

func main() {
	err := vfsgen.Generate(http.Dir("static"), vfsgen.Options{
		PackageName:  "rtprof",
		BuildTags:    "!dev",
		VariableName: "assets",
		Filename:     "assets_vfsdata.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
