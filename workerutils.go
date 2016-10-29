package main

import "fmt"

var (
	numWorkers = 5
)

// TaskMaster - overseer over the workers
type TaskMaster struct {
	Tasks        chan interface{}
	MaxWorkers   int
	Iter         int64
	StartedTasks int64
	Handler      handler
}

type handler interface{}

// InitTaskMaster - inititalize task master
func InitTaskMaster(numWorkers int, work handler) *TaskMaster {
	return &TaskMaster{MaxWorkers: numWorkers, Handler: work}
}

// AddTask - add new task to TaskMaster
func (tm *TaskMaster) AddTask(s interface{}) {
	tm.Tasks <- s
	tm.Iter++
	tm.StartedTasks++
}

func worker(id int, tasks chan string, quit <-chan bool) {
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			fmt.Printf("Worker %d Grab %s\n", id, task)
			Grab(task)
		case <-quit:
			return
		}
	}
}
