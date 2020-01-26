package main

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := newQueue()
	queue.put("1")
	if queue.len() != 1 {
		t.Errorf("Got %v queue len, want %v", queue.len(), 1)
	}
	if queue.cap() != 2 {
		t.Errorf("Got %v queue cap, want %v", queue.cap(), 2)
	}
	queue.put("2")
	if queue.len() != 2 {
		t.Errorf("Got %v queue len, want %v", queue.len(), 2)
	}
	if queue.cap() != 2 {
		t.Errorf("Got %v queue cap, want %v", queue.cap(), 2)
	}
	queue.put("3")
	if queue.len() != 3 {
		t.Errorf("Got %v queue len, want %v", queue.len(), 3)
	}
	if queue.cap() != 4 {
		t.Errorf("Got %v queue cap, want %v", queue.cap(), 4)
	}
	queue.put("4")
	queue.put("5")
	queue.put("6")
	queue.put("7")
	queue.put("8")
	queue.put("9")
	if queue.len() != 9 {
		t.Errorf("Got %v queue len, want %v", queue.len(), 9)
	}
	if queue.cap() != 16 {
		t.Errorf("Got %v queue cap, want %v", queue.cap(), 16)
	}
	value, ok := queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if value != "1" {
		t.Errorf("Got %v value, want %v", value, "1")
	}
	value, ok = queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if value != "2" {
		t.Errorf("Got %v value, want %v", value, "2")
	}
	value, ok = queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if value != "3" {
		t.Errorf("Got %v value, want %v", value, "3")
	}
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	value, _ = queue.get()
	if value != "9" {
		t.Errorf("Got %v value, want %v", value, "9")
	}
	_, ok = queue.get()
	if ok {
		t.Errorf("Got %v in queue get, want %v", ok, false)
	}
}

func BenchmarkQueue(b *testing.B) {
	queue := newQueue()
	b.ResetTimer()

	n := b.N
	for i := 0; i < n; i++ {
		queue.put("1")
	}
	for i := 0; i < n; i++ {
		value, ok := queue.get()
		if !ok {
			b.Errorf("Got %v error, want %v", ok, true)
		}
		if value != "1" {
			b.Errorf("Got %v value, want %v", value, "1")
		}
	}
}

func BenchmarkParallelQueue(b *testing.B) {
	queue := newQueue()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.put("1")
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, ok := queue.get()
			if !ok {
				b.Errorf("Got %v error, want %v", ok, true)
			}
			if value != "1" {
				b.Errorf("Got %v value, want %v", value, "1")
			}
		}
	})
}
