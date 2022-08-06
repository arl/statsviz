package example

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Work loops forever, generating a bunch of allocations of various sizes, plus
// some goroutines, in order to force the garbage collector to work.
func Work() {
	work(context.Background(), 0)
}

func work(ctx context.Context, lvl int) {
	m := make(map[int]interface{})
	tick := time.NewTicker(10 * time.Millisecond)

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			break
		}

		var obj interface{}

		switch rand.Intn(10) % 10 {
		case 0:
			obj = &struct {
				_ uint32
				_ uint16
			}{}
		case 1:
			obj = &struct {
				_ [3]uint64
			}{}
		case 2:
			obj = fmt.Sprintf("a relatively long and useless string %d", i)
		case 3:
			obj = make([]byte, i%1024)
		case 4, 5:
			obj = make([]byte, 10*i%1024)
		case 6, 7:
			obj = make([]string, 512)
			go func() {
				time.Sleep(time.Second)
			}()
		case 8:
			obj = make([]byte, 16*1024)
		case 9:
			if lvl > 2 {
				break
			}
			// Sometimes start another goroutine, just because...
			ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			work(ctx, lvl+1)
		}

		if i == 1000 {
			m = make(map[int]interface{})
			i = 0
		}

		m[i] = obj
	}
}
