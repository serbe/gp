package main

import "fmt"

var (
	numWorkers = 5
)

// TaskMaster - overseer over the workers
type TaskMaster struct {
	Tasks        chan interface{}
	Iter         int64
	StartedTasks int64
}

// InitTaskMaster - inititalize task master
func InitTaskMaster() TaskMaster {
	var tm TaskMaster

	return tm
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
			grab(task)
		case <-quit:
			return
		}
	}
}

func grab(host string) {
	body, err := fetch(host)
	if err != nil {
		finishTask <- true
		return
	}

	body = cleanBody(body)

	ips := getListIP(body)

	urls := getListURL(host, body)

	saveIP(ips)

	for _, item := range urls {
		if !urlList[item] {
			urlList[item] = true
			crawlChan <- item
		}
	}
	finishTask <- true
	return
}
