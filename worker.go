package main

import (
	"log"
	"sync"
	"time"
)

type worker struct {
	id     int64
	target string
	in     chan string
	out    chan Task
	quit   chan struct{}
	wg     *sync.WaitGroup
}

func (w *worker) start() {
	w.wg.Done()
	ticker := time.NewTicker(time.Duration(cfg.Timeout*3) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			log.Println("worker", w.id, "is sleep")
		case task := <-w.in:
			w.out <- crawl(task)
			ticker = time.NewTicker(time.Duration(cfg.Timeout*3) * time.Millisecond)
		case <-w.quit:
			w.wg.Done()
			return
		}
	}
}

func (w *worker) stop() {
	w.quit <- struct{}{}
}
