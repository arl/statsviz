# Easymon

Live monitoring of your app runtime stats (GC, MemStats, etc.) with a single import.

Similar to `import _ "net/http/pprof"`, import this package for its side-effects.


 - add `import _ "github.com/arl/easymon"` to your application
 - start an http server, if one isn't already running
 - points your browser to `http://host:port/debug/easymon`
 - enjoy
