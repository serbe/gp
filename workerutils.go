package main

var (
	numWorkers = 5
)

// TaskMaster - overseer over the workers
type TaskMaster struct {
	Tasks        chan interface{}
	Quit         chan bool
	MaxWorkers   int
	Iter         int64
	StartedTasks int64
	Handler      Handler
}

// Worker ...
type Worker struct {
	ID      int
	Work    chan interface{}
	Quit    chan bool
	Handler Handler
}

// Handler - any function
type Handler func(interface{})

// InitTaskMaster - inititalize task master
func InitTaskMaster(numWorkers int, work Handler) *TaskMaster {
	return &TaskMaster{MaxWorkers: numWorkers, Handler: work}
}

// AddTask - add new task to TaskMaster
func (tm *TaskMaster) AddTask(s interface{}) {
	tm.Tasks <- s
	tm.Iter++
	tm.StartedTasks++
}

// StartWorkers - start goroutines of nun workers
func (tm *TaskMaster) StartWorkers() {
	for i := 0; i < tm.MaxWorkers; i++ {
		// go worker(i, tm.Quit)
	}
}

// BeginWork - start loop for get and set channels
func (tm *TaskMaster) BeginWork() {
Loop:
	for {
		select {
		case newWork := <-crawlChan:
			tm.AddTask(newWork)
		case <-finishTask:
			// iter--
			// if iter == 0 {
			// 	for i := 0; i < numWorkers; i++ {
			// 		quit <- true
			// 	}
			break Loop
			// }
		}
	}
}

// Start worker
func (w Worker) Start() {
	go func() {
		for {
			select {
			case work := <-w.Work:
				w.Handler(work)
			case <-w.Quit:
				return
			}
		}
	}()
}
