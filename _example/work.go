package example

import (
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// Wwork loops forever, generating allocations of various sizes, in order to
// create artificial work for a nice 'demo effect'.
func Work() {
	m := make(map[int64]any)
	tick := time.NewTicker(30 * time.Millisecond)
	clearTick := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-clearTick.C:
			m = make(map[int64]any)
		case ts := <-tick.C:
			m[ts.UnixNano()] = newStruct()
		}
	}
}

// create a randomly sized struct (to create 'motion' on size classes plot).
func newStruct() any {
	nfields := rand.Intn(32)
	var fields []reflect.StructField
	for i := 0; i < nfields; i++ {
		fields = append(fields, reflect.StructField{
			Name:    "f" + strconv.Itoa(i),
			PkgPath: "main",
			Type:    reflect.TypeOf(""),
		})
	}
	return reflect.New(reflect.StructOf(fields)).Interface()
}
