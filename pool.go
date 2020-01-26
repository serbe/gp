package main

import (
	"time"

	"github.com/serbe/adb"
	"github.com/serbe/sites"
)

// Pool - specification of golang pool
type Pool struct {
	// running    bool
	nums    *nums
	input   Queue
	dp      *dbPool
	out     chan string
	quit    chan struct{}
	workers []worker
	cfg     *config
}

func newPool(cfg config) *Pool {
	nums := new(nums)
	db := adb.InitDB(cfg.DatabaseURL)
	dp := &dbPool{
		input: make(chan Task),
		nums:  nums,
		db:    &db,
		cfg:   &cfg,
	}
	p := &Pool{
		input: newQueue(),
		out:   make(chan string, cfg.Workers),
		dp:    dp,
		quit:  make(chan struct{}),
		nums:  nums,
		cfg:   &cfg,
	}
	p.startWorkers()
	go dp.start()
	go p.start()
	return p
}

func (p *Pool) startWorkers() {
	var (
		i       int64
		workers []worker
	)
	for i < p.cfg.Workers {
		worker := worker{
			id:     i,
			in:     p.out,
			out:    p.dp.input,
			quit:   make(chan struct{}),
			target: p.cfg.Target,
			nums:   p.nums,
		}
		go worker.start()
		workers = append(workers, worker)
		i++
	}
	p.workers = workers
}

func (p *Pool) start() {
	// p.running = true
	tick := time.Tick(time.Duration(200) * time.Microsecond)
	// ticker := time.NewTicker(time.Duration(p.cfg.Timeout*3) * time.Millisecond)
	for {
		select {
		// case <-ticker.C:
		// 	log.Println("Pool is sleep")
		// 	p.stop()
		case <-tick:
			if p.nums.getFreeWorkers() > 0 {
				value, ok := p.input.get()
				if ok {
					p.nums.decFreeWorkers()
					p.out <- value
					// ticker = time.NewTicker(time.Duration(p.cfg.Timeout*3) * time.Millisecond)
				}
			}
		case <-p.quit:
			p.dp.stop()
			for i := range p.workers {
				// p.wg.Add(1)
				p.workers[i].stop()
			}
			// close(p.quit)
			return
		}
	}
}

func (p *Pool) add(hostname string) error {
	if hostname == "" {
		return errEmptyHostname
	}
	// if !p.running {
	// 	return errNotRun
	// }
	// req := req{
	// 	ID:       p.addedTasks,
	// 	Hostname: hostname,
	// 	Proxy:    proxy,
	// }
	p.input.put(hostname)
	p.nums.incAddedTasks()
	// p.taskWG.Add(1)
	return nil
}

// func (p *Pool) stop() {
// 	p.quit <- struct{}{}
// 	// p.wg.Wait()
// 	// p.running = false
// }

// // IsRunning - check pool status is running
// func (p *Pool) IsRunning() bool {
// 	return p.running
// }

// // Wait - wait all task is done
// func (p *Pool) Wait() {
// 	p.taskWG.Wait()
// }

func (p *Pool) getHostList() {
	var list []string
	if p.cfg.UseFind {
		debugmsg("Start find proxy")
		newList := sites.ParseSites(p.cfg.LogDebug, p.cfg.LogErrors)
		saveToFile("newlist.txt", newList)
		oldList := p.dp.getLastProxy(100000)
		saveToFile("oldlist.txt", oldList)
		list = removeDuplicates(newList, oldList)
		saveToFile("list.txt", list)
		p.cfg.isUpdate = false
		debugmsg("End find proxy")
	} else if p.cfg.UseCheck {
		debugmsg("Start get proxy from db")
		if p.cfg.useTestLink {
			debugmsg("list empty for test link")
		} else if p.cfg.UseCheckAll {
			list = p.dp.getAllProxy()
		} else if p.cfg.UseFUP {
			list = p.dp.getFUPList()
		} else if p.cfg.UseCheckScheme {
			list = p.dp.getListWithScheme()
		} else {
			list = p.dp.getAllOld()
		}
		p.cfg.isUpdate = true
		debugmsg("End get proxy from db")
	}

	for i := range list {
		err := p.add(list[i])
		chkErr("add host to pool", err)
	}

	debugmsg("End get proxy from db")
}
