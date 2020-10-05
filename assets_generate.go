// +build ignore

package main

import (
	"encoding/json"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

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
	// Before generating the assets, generate the memstats.js file, it contains
	// the go documentation for the runtime.MemStats structure and is used inside
	// the web interface to show that documentation.
	jsbuf, err := genMemStatsDoc()
	if err != nil {
		log.Fatalln("extract memstats doc", err)
	}

	if err := ioutil.WriteFile("./static/memstats.js", jsbuf, 0644); err != nil {
		log.Fatalln("write memstats doc", err)
	}

	err = vfsgen.Generate(http.Dir("static"), vfsgen.Options{
		PackageName:  "statsviz",
		BuildTags:    "!dev",
		VariableName: "assets",
		Filename:     "assets_vfsdata.go",
	})
	if err != nil {
		log.Fatalln("generate assets", err)
	}
}

func genMemStatsDoc() ([]byte, error) {
	// Create the AST by parsing src and test.
	fset := token.NewFileSet()
	b, err := ioutil.ReadFile(filepath.Join(runtime.GOROOT(), "src", "runtime", "mstats.go"))
	if err != nil {
		return nil, err
	}

	f, err := parser.ParseFile(fset, "mstats.go", b, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	files := []*ast.File{f}

	// Compute package documentation.
	p, err := doc.NewFromFiles(fset, files, "runtime")
	if err != nil {
		return nil, err
	}

	js := make(map[string]string)

	tspec := p.Types[0].Decl.Specs[0].(*ast.TypeSpec)
	styp := tspec.Type.(*ast.StructType)
	for _, f := range styp.Fields.List {
		js[f.Names[0].Name] = f.Doc.Text()
	}

	return json.MarshalIndent(js, "", "  ")
}
