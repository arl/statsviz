package plot

import (
	"math"
	"runtime/metrics"
)

type goroutines struct {
	idxGoroutines int
	idxCreated    int
	idxNotInGo    int
	idxRunnable   int
	idxRunning    int
	idxWaiting    int

	lastCreated uint64
}

func makeGoroutines(indices ...int) metricsGetter {
	return &goroutines{
		idxGoroutines: indices[0],
		idxCreated:    indices[1],
		idxNotInGo:    indices[2],
		idxRunnable:   indices[3],
		idxRunning:    indices[4],
		idxWaiting:    indices[5],
		lastCreated:   math.MaxUint64,
	}
}

func (p *goroutines) values(samples []metrics.Sample) any {
	goroutines := samples[p.idxGoroutines].Value.Uint64()
	created := samples[p.idxCreated].Value.Uint64()
	notInGo := samples[p.idxNotInGo].Value.Uint64()
	runnable := samples[p.idxRunnable].Value.Uint64()
	running := samples[p.idxRunning].Value.Uint64()
	waiting := samples[p.idxWaiting].Value.Uint64()

	curCreated := created - p.lastCreated
	p.lastCreated = created

	return []uint64{
		goroutines,
		curCreated,
		notInGo,
		runnable,
		running,
		waiting,
	}
}

var goroutinesLayout = Scatter{
	Name:  "TODO(set later)",
	Title: "Goroutines",
	Type:  "scatter",
	Layout: ScatterLayout{
		Yaxis: ScatterYAxis{
			Title: "goroutines",
		},
	},
	Subplots: []Subplot{
		{
			Name:    "goroutines",
			Unitfmt: "%{y}",
		},
		{
			Name:    "created",
			Unitfmt: "%{y}",
			Type:    "bar",
		},
		{
			Name:    "not in Go",
			Unitfmt: "%{y}",
		},
		{
			Name:    "runnable",
			Unitfmt: "%{y}",
		},
		{
			Name:    "running",
			Unitfmt: "%{y}",
		},
		{
			Name:    "waiting",
			Unitfmt: "%{y}",
		},
	},
	InfoText: `<i>Goroutines</i> is <b>/sched/goroutines</b>, the count of live goroutines.
<i>Created</i> is the delta of <b>/sched/goroutines-created</b>, the cumulative number of created goroutines.
<i>Not in Go</i> is <b>/sched/goroutines/not-in-go</b>, the approximate count of goroutines running or blocked in a system call or cgo call.
<i>Runnable</i> is <b>/sched/goroutines/runnable</b>, the approximate count of goroutines ready to execute, but not executing.
<i>Running</i> is <b>/sched/goroutines/running</b>, the approximate count of goroutines executing.
<i>Waiting</i> is <b>/sched/goroutines/waiting</b>, the approximate count of goroutines waiting on a resource (I/O or sync primitives).`,
}
