package main

import (
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	w := worker{
		nums: new(nums),
		quit: make(chan struct{}),
		wg:   new(sync.WaitGroup),
	}

	w.run()
	if !w.running {
		t.Errorf("Got %v error, want %v", w.running, true)
	}
	if w.nums.getFreeWorkers() != 1 {
		t.Errorf("Got %v error, want %v", w.nums.getFreeWorkers(), 1)
	}
	w.wg.Add(1)
	w.stop()
	w.wg.Wait()
	if w.running {
		t.Errorf("Got %v error, want %v", w.running, false)
	}
	if w.nums.getFreeWorkers() != 0 {
		t.Errorf("Got %v error, want %v", w.nums.getFreeWorkers(), 0)
	}
}
