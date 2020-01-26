package main

import "sync"

type nums struct {
	sync.RWMutex
	addedTasks     int64
	completedTasks int64
	freeWorkers    uint64
	workProxy      int64
	anonProxy      int64
}

func (n *nums) getFreeWorkers() uint64 {
	n.RLock()
	defer n.RUnlock()
	return n.freeWorkers
}

func (n *nums) incFreeWorkers() {
	n.Lock()
	n.freeWorkers++
	n.Unlock()
}

func (n *nums) decFreeWorkers() {
	n.Lock()
	n.freeWorkers--
	n.Unlock()
}

func (n *nums) getAddedTasks() int64 {
	n.RLock()
	defer n.RUnlock()
	return n.addedTasks
}

func (n *nums) incAddedTasks() {
	n.Lock()
	n.addedTasks++
	n.Unlock()
}

func (n *nums) getCompletedTasks() int64 {
	n.RLock()
	defer n.RUnlock()
	return n.completedTasks
}

func (n *nums) incCompletedTasks() {
	n.Lock()
	n.completedTasks++
	n.Unlock()
}

// func (n *nums) decCompletedTasks() {
// 	n.Lock()
// 	n.completedTasks--
// 	n.Unlock()
// }
