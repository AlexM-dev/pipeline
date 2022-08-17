package ring

import (
	"fmt"
	"sync"
)

type IntRing struct {
	data []int
	pos  int
	size int
	m    sync.Mutex
}

func NewRing(size int) *IntRing {
	return &IntRing{make([]int, size), -1, size, sync.Mutex{}}
}

func (r *IntRing) Write(n int) {
	r.m.Lock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.data[i-1] = r.data[i]
		}
		r.data[r.pos] = n
	} else {
		r.pos++
		r.data[r.pos] = n
	}
	r.m.Unlock()
}

func (r *IntRing) ReadAll() []int {
	if r.pos < 0 {
		return nil
	}
	r.m.Lock()
	var output = r.data[:r.pos+1]
	r.pos = -1
	r.m.Unlock()
	return output
}

func (r IntRing) Size() int {
	return r.size
}

func (r *IntRing) RemoveByIndex(index int) error {
	if index < 0 || index >= r.size {
		return fmt.Errorf("несуществующий индекс %d", index)
	}
	if r.data[index] != 0 {
		r.data[index] = 0
	}
	return nil
}
