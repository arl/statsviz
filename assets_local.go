// +build dev

package rtprof

import "net/http"

// assets contains project assets located in current directory.
var assets http.FileSystem = http.Dir("static")
