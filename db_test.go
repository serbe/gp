package main

import "testing"

import "sync"

func newTestDP() dbPool {
	dp := dbPool{
		input: make(chan Task),
		quit:  make(chan struct{}),
		// nums    *nums
		// db      *adb.DB
		// cfg     *config
		wg: new(sync.WaitGroup),
	}
	return dp
}

func TestDP(t *testing.T) {
	dp := newTestDP()
	dp.run()
	if !dp.running {
		t.Errorf("Got %v error, want %v", dp.running, true)
	}
	dp.wg.Add(1)
	dp.stop()
	dp.wg.Wait()
	if dp.running {
		t.Errorf("Got %v error, want %v", dp.running, false)
	}
}
