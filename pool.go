package main

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errEmptyTarget = errors.New("error: empty target hostname")
	errNotRun      = errors.New("error: pool is not running")
	errNotWait     = errors.New("error: pool is not waiting tasks")
)

// Pool - specification of golang pool
type Pool struct {
	running        bool
	useOutChan     bool
	numWorkers     int64
	addedTasks     int64
	completedTasks int64
	out            Queue
	toWork         chan string
	fromWork       chan Task
	outTasks       chan Task
	quit           chan struct{}
	workers        []worker
	wg             sync.WaitGroup
	taskWG         sync.WaitGroup
}

// NewPool - create new goroutine pool with channels
// numWorkers - max workers
func NewPool(numWorkers int64) *Pool {
	// rand.Seed(time.Now().UnixNano())
	p := &Pool{
		numWorkers: numWorkers,
		out:        newQueue(),
		toWork:     make(chan string, numWorkers),
		fromWork:   make(chan Task, numWorkers),
		outTasks:   make(chan Task, 1),
		quit:       make(chan struct{}),
	}
	p.wg.Add(1)
	p.startWorkers()
	go p.start()
	p.wg.Wait()
	return p
}

func (p *Pool) startWorkers() {
	var (
		i       int64
		workers []worker
	)
	for i < p.numWorkers {
		p.wg.Add(1)
		worker := worker{
			id:   i,
			in:   p.toWork,
			out:  p.fromWork,
			quit: make(chan struct{}),
			wg:   &p.wg,
		}
		go worker.start()
		workers = append(workers, worker)
		i++
	}
	p.workers = workers
}

func (p *Pool) start() {
	p.running = true
	p.wg.Done()
	// tick := time.Tick(time.Duration(200) * time.Microsecond)
	ticker := time.NewTicker(time.Duration(cfg.Timeout*3) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			log.Println("Pool is sleep")
		// case <-tick:
		// 	task, ok := p.in.get()
		// 	if ok {
		// 		p.toWork <- task
		// 	}
		case task := <-p.fromWork:
			p.outTasks <- task
			p.incCompleted()
			ticker = time.NewTicker(time.Duration(cfg.Timeout*3) * time.Millisecond)
			p.taskWG.Done()
		case <-p.quit:
			for i := range p.workers {
				p.wg.Add(1)
				p.workers[i].stop()
			}
			// close(p.quit)
			return
		}
	}
}

// Add - adding task to pool
func (p *Pool) Add(hostname string, proxy string) error {
	if hostname == "" {
		return errEmptyTarget
	}
	if !p.running {
		return errNotRun
	}
	// req := req{
	// 	ID:       p.addedTasks,
	// 	Hostname: hostname,
	// 	Proxy:    proxy,
	// }
	p.incAdded()
	// p.in.put(req)
	p.taskWG.Add(1)
	return nil
}

// Stop - stop pool and all workers
func (p *Pool) Stop() {
	p.quit <- struct{}{}
	p.wg.Wait()
	p.running = false
}

// IsRunning - check pool status is running
func (p *Pool) IsRunning() bool {
	return p.running
}

// Wait - wait all task is done
func (p *Pool) Wait() {
	p.taskWG.Wait()
}

func (p *Pool) incAdded() {
	atomic.AddInt64(&p.addedTasks, 1)
}

// Added - number of adding tasks
func (p *Pool) Added() int64 {
	return atomic.LoadInt64(&p.addedTasks)
}

func (p *Pool) incCompleted() {
	atomic.AddInt64(&p.completedTasks, 1)
}

// Completed - number of completed tasks
func (p *Pool) Completed() int64 {
	return atomic.LoadInt64(&p.completedTasks)
}
