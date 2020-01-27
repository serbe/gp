package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// // import (
// // 	"fmt"
// // 	"net/http"
// // 	"net/http/httptest"
// // 	"testing"
// // 	"time"
// // )

// var (
// 	numWorkers int64 = 4
// 	t10ms            = time.Duration(10) * time.Millisecond
// 	t30ms            = time.Duration(30) * time.Millisecond
// 	testList         = []string{
// 		"http://41.217.219.49*:51283",
// 		"http://197.232.69.137*:33053",
// 		"http://68.183.33.150:3128",
// 		"http://171.7.69.234:8213",
// 		"http://36.66.203.127:8080",
// 	}
// )

// func testHandler(w http.ResponseWriter, _ *http.Request) {
// 	fmt.Fprint(w, "Test page")
// }

// func testHandlerWithTimeout(w http.ResponseWriter, _ *http.Request) {
// 	time.Sleep(t30ms)
// 	fmt.Fprint(w, "Test page with timeout")
// }

// func TestClosedInputTaskChanByTimeout(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()

// 	p := NewPool(numWorkers)
// 	// p.SetQuitTimeout(10)
// 	err := p.Add(ts.URL, "")
// 	if err != nil {
// 		t.Errorf("Got %v error, want %v", err, nil)
// 	}
// 	time.Sleep(t30ms)
// 	err = p.Add(ts.URL, "")
// 	if err == nil {
// 		t.Errorf("Got %v error, want %v", nil, errNotRun)
// 	}
// }

// func TestNoServer(t *testing.T) {
// 	cfg := config{
// 		Workers: numWorkers,
// 	}
// 	p := newPool(cfg)
// 	p.run()
// 	if !p.running {
// 		t.Errorf("pool is %v, want %v", p.running, true)
// 	}
// 	if int64(len(p.workers)) != numWorkers {
// 		t.Errorf("pool have %v numWorkers, want %v", len(p.workers), numWorkers)
// 	}
// 	if len(p.workers) != int(numWorkers) {
// 		t.Errorf("pool have %v Workers, want %v", len(p.workers), numWorkers)
// 	}
// 	err := p.add("")
// 	if err != errEmptyHostname {
// 		t.Errorf("Got %v error, want %v", err, errEmptyHostname)
// 	}
// 	if p.nums.getAddedTasks() != 0 {
// 		t.Errorf("Wrong input added tasks, got %v, want %v", p.nums.getAddedTasks(), 0)
// 	}
// 	err = p.add(":")
// 	if err != nil {
// 		t.Errorf("Got %v error, want %v", err, nil)
// 	}
// 	err = p.add("http://127.0.0.1:80")
// 	if err != nil {
// 		t.Errorf("Got %v error, want %v", err, nil)
// 	}
// 	err = p.add("http://127.0.0.1:80")
// 	if err != nil {
// 		t.Errorf("Got %v error, want %v", err, nil)
// 	}
// 	if p.nums.getAddedTasks() != 3 {
// 		t.Errorf("Wrong input jobs. got %v, want %v", p.nums.getAddedTasks(), 3)
// 	}
// 	p.stop()
// 	if p.running {
// 		t.Errorf("pool is %v, want %v", p.running, false)
// 	}
// 	// err = p.add("http://127.0.0.1:80/", "")
// 	// if err != errNotRun {
// 	// 	t.Errorf("Got %v error, want %v", err, errNotRun)
// 	// }
// }

// func TestWithServer(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()

// 	p := NewPool(numWorkers)
// 	err := p.Add(ts.URL, "")
// 	if err != nil {
// 		t.Errorf("Got %v with adding task, want %v", err, nil)
// 	}
// 	p.taskWG.Wait()
// 	if p.completedTasks != p.addedTasks {
// 		t.Errorf("Got %v completed tasks, want %v", p.completedTasks, p.addedTasks)
// 	}
// 	task, ok := p.Get()
// 	if !ok {
// 		t.Errorf("Got %v with getting task, want %v", ok, true)
// 	}
// 	if string(task.Body) != "Test page" {
// 		t.Errorf("Got %v in task.Body, want '%v'", string(task.Body), "Test page")
// 	}
// 	if p.completedTasks != 1 {
// 		t.Errorf("Got %v completed tasks, want %v", p.completedTasks, 1)
// 	}
// 	err = p.Add(ts.URL, "")
// 	if err != nil {
// 		t.Errorf("Got %v with adding task, want %v", err, nil)
// 	}
// 	p.Wait()
// 	_, ok = p.Get()
// 	if !ok {
// 		t.Errorf("Got %v with getting task, want %v", ok, true)
// 	}
// 	if p.completedTasks != 2 {
// 		t.Errorf("Got %v error, want %v", p.completedTasks, 2)
// 	}
// 	if p.completedTasks != p.Completed() {
// 		t.Errorf("Got %v error, want %v", p.completedTasks, p.Completed())
// 	}
// 	p.Stop()
// }

// func TestWithTimeout(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandlerWithTimeout))
// 	defer ts.Close()

// 	p := NewPool(numWorkers)
// 	p.NetTimeout(100)
// 	if timeout != 100 {
// 		t.Errorf("Got %v net timeout, want %v", timeout, 100)
// 	}
// 	_ = p.Add(ts.URL, "")
// 	p.Wait()
// 	task, ok := p.Get()
// 	if !ok {
// 		t.Errorf("Got %v with getting task, want %v", ok, true)
// 	}
// 	if string(task.Body) != "Test page with timeout" {
// 		t.Errorf("Got %v error, want '%v'", string(task.Body), "Test page with timeout")
// 	}
// 	p.NetTimeout(5)
// 	if timeout != 5 {
// 		t.Errorf("Got %v net timeout, want %v", timeout, 5)
// 	}
// 	_ = p.Add(ts.URL, "")
// 	p.Wait()
// 	task, _ = p.Get()
// 	if task.Error == nil {
// 		t.Errorf("Got no error, want %v", task.Error)
// 	}
// 	if p.completedTasks != 2 {
// 		t.Errorf("Got %v error, want %v", p.completedTasks, 2)
// 	}
// 	p.Stop()
// }

// func TestOutChan(t *testing.T) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()

// 	p := NewPool(numWorkers)
// 	if p.IsUseOutChan() {
// 		t.Errorf("Got %v on use out chan, want %v", p.IsUseOutChan(), false)
// 	}
// 	ch := p.UseOutChan()
// 	if !p.IsUseOutChan() {
// 		t.Errorf("Got %v on use out chan, want %v", p.IsUseOutChan(), true)
// 	}
// 	_ = p.Add(ts.URL, "")
// 	task := <-ch
// 	if string(task.Body) != "Test page" {
// 		t.Errorf("Got %v in task.Body, want '%v'", string(task.Body), "Test page")
// 	}
// 	if p.Completed() != 1 {
// 		t.Errorf("Got %v completed tasks, want %v", p.Completed(), 1)
// 	}
// 	p.Stop()
// }

// // func TestWaitingTasks(t *testing.T) {
// // 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// // 	defer ts.Close()

// // 	p := NewPool(1)
// // 	_ = p.Add(ts.URL, "")
// // 	_ = p.Add(ts.URL, "")
// // 	p.EndWaitingTasks()
// // 	for range p.ResultChan {
// // 	}
// // 	if p.GetCompletedTasks() != 2 {
// // 		t.Errorf("Got %v error, want %v", p.GetCompletedTasks(), 1)
// // 	}
// // }

// func BenchmarkAccumulate(b *testing.B) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()
// 	// b.ResetTimer()

// 	p := NewPool(numWorkers)
// 	n := b.N
// 	for i := 0; i < n; i++ {
// 		err := p.Add(ts.URL, "")
// 		if err != nil {
// 			b.Errorf("Got %v error, want %v", err, nil)
// 		}
// 	}
// }

// func BenchmarkFullProcess(b *testing.B) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()
// 	// b.ResetTimer()

// 	p := NewPool(numWorkers)
// 	n := b.N
// 	for i := 0; i < n; i++ {
// 		err := p.Add(ts.URL, "")
// 		if err != nil {
// 			b.Errorf("Got %v error, want %v", err, nil)
// 		}
// 	}
// 	// p.Wait()
// 	for i := 0; i < n; i++ {
// 		_, _ = p.Get()
// 		// if !ok {
// 		// 	b.Errorf("Got %v with getting task, want %v", ok, true)
// 		// }
// 		// if task.Error != nil {
// 		// 	b.Errorf("Task %v have error %v", task.ID, task.Error)
// 		// }
// 	}
// }

// // func BenchmarkParallel(b *testing.B) {
// // 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// // 	defer ts.Close()
// // 	// b.ResetTimer()

// // 	p := NewPool(numWorkers)
// // 	b.RunParallel(func(pb *testing.PB) {
// // 		for pb.Next() {
// // 			err := p.Add(ts.URL, "")
// // 			if err != nil {
// // 				b.Errorf("Got %v error, want %v", err, nil)
// // 			}
// // 		}
// // 	})
// // 	b.RunParallel(func(pb *testing.PB) {
// // 		for pb.Next() {
// // 			_, _ = p.Get()
// // 		}
// // 	})
// // }

// func BenchmarkOutChan(b *testing.B) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()
// 	// b.ResetTimer()

// 	p := NewPool(numWorkers)
// 	ch := p.UseOutChan()
// 	n := b.N
// 	for i := 0; i < n; i++ {
// 		err := p.Add(ts.URL, "")
// 		if err != nil {
// 			b.Errorf("Got %v error, want %v", err, nil)
// 		}
// 	}
// 	for i := 0; i < n; i++ {
// 		task, ok := <-ch
// 		if !ok {
// 			b.Errorf("Got %v with getting task, want %v", ok, true)
// 		}
// 		if task.Error != nil {
// 			b.Errorf("Task %v have error %v", task.ID, task.Error)
// 		}
// 	}
// }

// func BenchmarkParallelOutChan(b *testing.B) {
// 	ts := httptest.NewServer(http.HandlerFunc(testHandler))
// 	defer ts.Close()
// 	// b.ResetTimer()

// 	p := NewPool(numWorkers)
// 	ch := p.UseOutChan()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			err := p.Add(ts.URL, "")
// 			if err != nil {
// 				b.Errorf("Got %v error, want %v", err, nil)
// 			}
// 		}
// 	})
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			select {
// 			case task, ok := <-ch:
// 				if !ok {
// 					b.Errorf("Got %v with getting task, want %v", ok, true)
// 				}
// 				if task.Error != nil {
// 					b.Errorf("Task %v have error %v", task.ID, task.Error)
// 				}
// 			}
// 		}
// 	})
// }
