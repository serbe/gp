package main

import "testing"

func TestDP(t *testing.T) {

	dp := dbPool{
		input: make(chan Task),
		quit:  make(chan struct{}),
		// nums    *nums
		// db      *adb.DB
		// cfg     *config
	}

	dp.start()
	if !dp.running {
		t.Errorf("Got %v error, want %v", dp.running, true)
	}
}
