package main

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/serbe/adb"

)

type dbPool struct {
	running bool
	input   chan Task
	quit    chan struct{}
	nums    *nums
	db      *adb.DB
	cfg     *config
}

func (dp *dbPool) start() {
	dp.running = true
	for {
		select {
		case task := <-dp.input:
			proxy := taskToProxy(task, dp.cfg)

			if dp.cfg.isUpdate {
				chkErr("saveProxy Update "+proxy.Hostname, dp.db.Update(&proxy))
			} else {
				chkErr("saveProxy Insert "+proxy.Hostname, dp.db.Insert(&proxy))
			}
			if proxy.IsWork {
				dp.nums.workProxy++
				if proxy.IsAnon {
					dp.nums.anonProxy++
				}
			}
			dp.nums.incCompletedTasks()
			debugmsg(fmt.Sprintf("%d/%d/%d %v %v %v",
				dp.nums.getAddedTasks(),
				dp.nums.getCompletedTasks(),
				dp.nums.workProxy,
				proxy.IsAnon,
				proxy.Hostname,
				task.Error,
			))

		case <-dp.quit:
			debugmsg("dbPool quit")
			dp.running = false
			return
		}
	}
}

func (dp *dbPool) stop() {
	dp.quit <- struct{}{}
}

func (dp *dbPool) getLastProxy(num int64) []string {
	list, err := dp.db.GetLast(num)
	chkErr("lastUpdated", err)
	return list
}

func (dp *dbPool) getAllProxy() []string {
	list, err := dp.db.GetAll()
	chkErr("getAllProxy", err)
	return list
}

func (dp *dbPool) getFUPList() []string {
	var list []string
	hosts, err := dp.db.GetUniqueHosts()
	chkErr("getFUPList ProxyGetUniqueHosts", err)
	ports, err := dp.db.GetFrequentlyUsedPorts()
	chkErr("getFUPList ProxyGetFrequentlyUsedPorts", err)
	for _, host := range hosts {
		for _, port := range ports {
			u := "http://" + host + ":" + strconv.Itoa(port)
			_, err := url.Parse(u)
			if err == nil {
				list = append(list, u)
			}
		}
	}
	return list
}

func (dp *dbPool) getListWithScheme() []string {
	var newList []string
	list, err := dp.db.GetAllScheme(HTTP)
	chkErr("getListWithScheme ProxyGetAllScheme", err)
	for i := range list {
		u, err := url.Parse(list[i])
		if err == nil {
			newList = append(newList, "https://"+u.Host+":"+u.Port())
			newList = append(newList, "socks5://"+u.Host+":"+u.Port())
		}
	}
	return newList
}

func (dp *dbPool) getAllOld() []string {
	list, err := dp.db.GetAllOld()
	chkErr("getProxyListFromDB ProxyGetAllOld", err)
	return list
}
