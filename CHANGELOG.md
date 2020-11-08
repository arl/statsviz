v0.2.1 / 2020-10-29
===================

  * Fix websocket handler now working with https [#25](https://github.com/arl/statsviz/pull/25)

v0.2.0 / 2020-10-25
===================

  * `Register` now accepts options (functional options API) [#20](https://github.com/arl/statsviz/pull/20):
    + `Root` allows to root statsviz at a path different than `/debug/statsviz`
    + `SendFrequency` allows to set the frequency at which stats are emitted.

v0.1.1 / 2020-10-12
===================

  * Do not leak timer in sendStats

v0.1.0 / 2020-10-10
===================

  * First released version
