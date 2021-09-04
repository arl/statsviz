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

The user interface aims to be simple, light and minimal.

Assets are located in the `internal/static` directory and are embedded with
[`go:embed`](https://pkg.go.dev/embed).


## Contribute by improving documentation

No contribution is too small, improvements to code comments and/or README
are welcome!

Thank you!


## Contribute by adding an example

There are many Go libraries to handle HTTP routing.

Feel free to add an example to show how to register statsviz with your favourite
library.

Please add a directory under `./_example`. For instance, if you want to add an
example showing how to register statsviz within library `foobar`:

 - create a directory `./_example/foobar/`
 - create a file `./_example/foobar/main.go`
 - call `go example.Work()` as the first line of your example (see other
   examples). This forces the garbage collector to _do something_ so that
   statsviz interface won't remain static when an user runs your example.
 - the code should be `gofmt`ed
 - the example should compile and run
 - when ran, statsviz interface should be accessible at http://localhost:8080/debug/statsviz

Thanks a lot!
