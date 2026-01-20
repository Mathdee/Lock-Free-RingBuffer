package ringbuffer

import (
	"fmt"
	"sync/atomic"
)

//Quick Notes:
// LAMPORT'S SPSC: CONSUMER OWNS TAIL, PRODUCER OWNS HEAD
// Head and Tail are placed 64-bytes apart so they live on different cache lines (typically L1 cache line size is 64 bytes)
// In Go, `sync/atomic` operations are sequentially consistent by default, so no need for explicit memory barriers like mutexes)
//
// Empty when --> head == tail
// Full when --> head - tail == size
//

type RingBuffer struct {
	// head is written by producer only, read by consumer
	head  uint64   // 8 bytes per cache line
	_pad0 [56]byte // 56 bytes,  head+pad = 64 bytes cache line

	// tail is written by consumer only, read by producer
	tail  uint64   // 8 bytes per cache line
	_pad1 [56]byte // 56-bytes padding so tail is on its own 64-bytes cache line

	size uint64
	mask uint64
	buf  []uint64
}

func NewRingBuffer(size uint64) *RingBuffer {

	if size < 2 || (size&(size-1) != 0) {
		panic(fmt.Sprintf("Careful - size must be a power of 2: current size is %d", size))
	}

	return &RingBuffer{
		size: size,
		mask: size - 1,
		buf:  make([]uint64, size),
	}
}

// This method adds a value to the ring buffer.
// Step 1: Load the current head and tail values.
// Step 2: Check if the ring buffer is full.
// Step 3: Add the value to the ring buffer.
// Step 4: Increment the head value.
// Step 5: Return true if the value was added successfully.
// Step 6: Return false if the value was not added successfully.
func (rb *RingBuffer) Push(val uint64) bool {
	h := atomic.LoadUint64(&rb.head)
	t := atomic.LoadUint64(&rb.tail)

	if h-t == rb.size {
		return false
	}

	rb.buf[h&rb.mask] = val
	atomic.StoreUint64(&rb.head, h+1)
	return true
}

// This method removes a value from the ring buffer.
// Step 1: Load the current head and tail values.
// Step 2: Check if the ring buffer is empty.
// Step 3: Remove the value from the ring buffer.
// Step 4: Increment the tail value.
// Step 5: Return the value and true if the value was removed successfully.
// Step 6: Return 0 and false if the value was not removed successfully.
func (rb *RingBuffer) Pop() (uint64, bool) {
	h := atomic.LoadUint64(&rb.head)
	t := atomic.LoadUint64(&rb.tail)

	if t == h {
		return 0, false
	}

	val := rb.buf[t&rb.mask]
	atomic.StoreUint64(&rb.tail, t+1)
	return val, true

}
