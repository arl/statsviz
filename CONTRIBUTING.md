Contributing
============

First of all, thank you to consider contributing to this open-source project!

Pull-requests are welcome!


## Contribute to statviz Go library

The statsviz Go API is very thin so there's not much to do and it's unlikely that
the API will change, however some new options can be added to `statsviz.Register` 
without breaking compatibility.
That being said, there may be things to improve in the implementation, any
contribution is very welcome!

If you've decided to contribute, thank you so much, please comment on the existing 
issue or create one stating you want to tackle it, so we can assign it to you and 
reduce the possibility of duplicate work.


## Contribute to the user interface (html/css/javascript)

The user interface aims to be light, minimal, simple.

This program uses [vfsgen](github.com/shurcooL/vfsgen) to embed the content of 
the `/static` directory into the final binary. `vfsgen` generates the `assets`
variable in `assets_vfsdata.go`. `assets` statically implements an 
`http.FileSystem` rooted at `/static/` which contains the files statsviz serves.

While working on statsviz web interface, it's easier to directly serve the
content of the `/static` directory than regenerating the assets after each 
modification. Passing `-tags dev` to `go build` will do just that, the
directory served will be that of your filesystem.

To commit some changes of the files in the `/static` directory, `assets`
must be regenerated (or the CI will complain anyway).
To do so just call, from the project root:

```
go generate
go mod tidy
```

With Go modules enabled, this will download the latest version of 
`github.com/shurcooL/vfsgen` and update `assets_vfsdata.go` so that 
it reflects the new content of the `/static` directory. Then, 
commits both the changes to `/static` and those to `assets_vfsdata.go`.
