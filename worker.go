package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type worker struct {
	id      int64
	running bool
	timeout int64
	target  string
	in      chan string
	out     chan Task
	quit    chan struct{}
	nums    *nums
	wg      *sync.WaitGroup
}

func (w *worker) run() {
	w.wg.Add(1)
	go w.start()
	w.wg.Wait()
}

func (w *worker) start() {
	w.nums.incFreeWorkers()
	w.running = true
	w.wg.Done()
	for {
		select {
		case hostname := <-w.in:
			w.out <- w.crawl(hostname)
			w.nums.incFreeWorkers()
		case <-w.quit:
			w.running = false
			w.nums.decFreeWorkers()
			w.wg.Done()
			return
		}
	}
}

func (w *worker) stop() {
	w.quit <- struct{}{}
}

func (w *worker) crawl(proxyURL string) Task {
	startTime := time.Now()
	var (
		proxy *url.URL
		task  Task
		err   error
	)
	// log.Println(w.target, proxyURL)
	task.Proxy = proxyURL
	if task.Proxy != "" {
		proxy, err = url.Parse(task.Proxy)
		if err != nil {
			task.Error = err
			return task
		}
	}
	client := &http.Client{
		Timeout: time.Duration(w.timeout) * time.Millisecond,
	}
	client.Transport = &http.Transport{
		Proxy:             http.ProxyURL(proxy),
		DisableKeepAlives: true,
	}
	request, err := http.NewRequest(http.MethodGet, w.target, nil)
	if err != nil {
		task.Error = err
		return task
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	request.Header.Set("Connection", "close")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Referer", "https://www.google.com/")
	response, err := client.Do(request)
	if err != nil {
		if response != nil {
			_ = response.Body.Close()
		}
		task.Error = err
		return task
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		task.Error = err
		err = response.Body.Close()
		if err != nil {
			task.Error = err
		}
		return task
	}
	task.Body = body
	task.Response = time.Since(startTime)
	err = response.Body.Close()
	if err != nil {
		task.Error = err
	}
	return task
}
