package main

import (
	"regexp"
	"sync"
	"time"

	"github.com/serbe/adb"
	"github.com/serbe/pool"
)

type mapLink struct {
	sync.RWMutex
	values map[string]adb.Link
}

func newMapLink() *mapLink {
	return &mapLink{values: make(map[string]adb.Link)}
}

func (ml *mapLink) fillMapLink(linkList []adb.Link) {
	for _, link := range linkList {
		ml.set(link)
	}
}

func (ml *mapLink) get(hostname string) (adb.Link, bool) {
	ml.RLock()
	link, ok := ml.values[hostname]
	ml.RUnlock()
	return link, ok
}

func (ml *mapLink) set(link adb.Link) {
	ml.Lock()
	ml.values[link.Hostname] = link
	ml.Unlock()
}

func (ml *mapLink) update(hostname string) {
	ml.Lock()
	link := ml.values[hostname]
	link.Update = true
	link.UpdateAt = time.Now()
	ml.values[hostname] = link
	ml.Unlock()
}

func newLink(hostname string, insert bool, iterate bool) adb.Link {
	var link adb.Link
	link.Hostname = hostname
	link.Insert = insert
	link.Iterate = iterate
	return link
}

func (ml *mapLink) existLink(hostname string) bool {
	_, ok := ml.get(hostname)
	return ok
}

func getMapLink() *mapLink {
	ml := newMapLink()
	if useTestLink {
		link := newLink(testLink, false, true)
		ml.set(link)
	} else if useAddLink {
		link := newLink(addLink, true, true)
		ml.set(link)
	} else {
		ml.fillMapLink(db.LinksGetAll())
	}
	return ml
}

func (ml *mapLink) saveAll() {
	debugmsg("start saveAllLinks")
	var (
		u, i int64
	)
	for _, l := range ml.values {
		if l.Insert {
			i++
			db.LinkInsert(l)
		}
		if l.Update {
			u++
			db.LinkUpdate(l)
		}
	}
	debugmsg("update links", u)
	debugmsg("insert links", i)
	debugmsg("end saveAllLinks")
}

func (ml *mapLink) getNewLinksFromTask(task *pool.Task) []adb.Link {
	var links []adb.Link
	for i := range reURL {
		host, err := getHost(task.Hostname)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(reURL[i])
		if !re.Match(task.Body) {
			continue
		}
		allResults := re.FindAllSubmatch(task.Body, -1)
		for _, result := range allResults {
			hostname := host + "/" + string(result[1])
			if ml.existLink(hostname) {
				continue
			}
			link := newLink(hostname, true, false)
			link.UpdateAt = time.Now()
			ml.set(link)
			links = append(links, link)
		}
	}
	return links
}
