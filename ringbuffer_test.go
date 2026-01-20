package ringbuffer

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkRingBuffer(b *testing.B) {
	b.ReportAllocs()

	old := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(old)

	rb := NewRingBuffer(1 << 16)

	var start uint32
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for atomic.LoadUint32(&start) == 0 {
			runtime.Gosched()
		}
		for i := 0; i < b.N; i++ {
			for {
				if _, ok := rb.Pop(); ok {
					break
				}
				runtime.Gosched()

			}
		}
	}()

	b.ResetTimer()
	atomic.StoreUint32(&start, 1)

	for i := 0; i < b.N; i++ {
		v := uint64(i)
		for {
			if rb.Push(v) {
				break
			}
			runtime.Gosched()
		}

	}
	wg.Wait()
}

func BenchmarkRingChan(b *testing.B) {
	b.ReportAllocs()
	old := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(old)

	ch := make(chan uint64, 1<<16)

	var start uint32
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for atomic.LoadUint32(&start) == 0 {
			runtime.Gosched()
		}
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()

	b.ResetTimer()

	atomic.StoreUint32(&start, 1)

	for i := 0; i < b.N; i++ {
		ch <- uint64(i)
	}

	wg.Wait()

}
