package example

import (
	"fmt"
	"time"
)

// Work loops forever, generating a bunch of allocations of various sizes in
// order to force the garbage collector to work.
func Work() {
	m := make(map[int]interface{})
	i := 0
	for ; ; i++ {

		var obj interface{}
		switch i % 6 {
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
			obj = fmt.Sprint("a relatively long and useless string %d", i)
		case 3:
			obj = make([]byte, i%1024)
		case 4:
			obj = make([]byte, 10*i%1024)
		case 5:
			obj = make([]string, 512)
		}

		if i == 1000 {
			m = make(map[int]interface{})
			i = 0
		}

		m[i] = obj
		time.Sleep(10 * time.Millisecond)
	}
}
