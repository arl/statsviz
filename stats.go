package statsviz

import "runtime"

type stats struct {
	Mem          runtime.MemStats
	NumGoroutine int
}
