package imaging

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var parallelizationEnabled = true

// if GOMAXPROCS = 1: no goroutines used
// if GOMAXPROCS > 1: spawn N=GOMAXPROCS workers in separate goroutines
func parallel(dataSize int, fn func(partStart, partEnd int)) {
	numGoroutines := 1
	partSize := dataSize

	if parallelizationEnabled {
		numProcs := runtime.GOMAXPROCS(0)
		if numProcs > 1 {
			numGoroutines = numProcs
			partSize = dataSize / (numGoroutines * 100)
			if partSize < 1 {
				partSize = 1
			}
		}
	}

	if numGoroutines == 1 {
		fn(0, dataSize)
	} else {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		idx := uint64(0)

		for p := 0; p < numGoroutines; p++ {
			go func() {
				defer wg.Done()
				for {
					partStart := int(atomic.AddUint64(&idx, uint64(partSize))) - partSize
					if partStart >= dataSize {
						break
					}
					partEnd := partStart + partSize
					if partEnd > dataSize {
						partEnd = dataSize
					}
					fn(partStart, partEnd)
				}
			}()
		}

		wg.Wait()
	}
}
