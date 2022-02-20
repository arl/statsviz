Unreleased yet
==============
  * Switch javascript code to ES6 (#65)
  * Build and test all examples (#63)
  * Assets are `go:embed`ed, so the minimum go version is now go1.16 (#55)
  * Polishing (README, small UI improvements) (#54)
  * Small ui improvements: link to go.dev rather than golang.org

v0.4.0 / 2021-05-08
==================

  * Auto-reconnect to new server from GUI after closed websocket connection (#49)
  * Reorganize examples (#51)
  * Make `IndexAtRoot` returns an `http.HandlerFunc` instead of `http.Handler` (#52)

v0.3.0 / 2021-02-14
==================

  * Enable 'save as png' button on plots (#44)

v0.2.2 / 2020-12-13
==================

  * Use Go Modules for 'github.com/gorilla/websocket' (#39)
  * Support custom frequency (#37)
  * Added fixed go-chi example (#38)
  * `_example`: add echo (#22)
  * `_example`: add gin example (#34)
  * ci: track coverage
  * RegisterDefault returns an error now
  * Ensure send frequency is a strictly positive integer
  * Don't log if we can't upgrade to websocket
  * `_example`_example: add chi router (#38)
  * `_example`_example: change structure to have one example per directory
v0.2.1 / 2020-10-29
===================

  * Fix websocket handler now working with https (#25)

v0.2.0 / 2020-10-25
===================

  * `Register` now accepts options (functional options API) (#20)
    + `Root` allows to root statsviz at a path different than `/debug/statsviz`
    + `SendFrequency` allows to set the frequency at which stats are emitted.

v0.1.1 / 2020-10-12
===================

  * Do not leak timer in sendStats

v0.1.0 / 2020-10-10
===================

  * First released version
