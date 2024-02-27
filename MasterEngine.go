package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "runtime"
	// "sync/atomic"
)

// var engineDriver = FangTianXiaEngineDriver{
// rootLinkSpec : RootLinkSpec{url: "https://sh.esf.fang.com/newsecond/esfcities.aspx",
// 	linkStartFlag: "<div class=\"outCont\" id=\"c02\"",
// 	linkEndFlag:   "</div>"},
// }

var engineDriver = LianJiaEngineDriver{
	rootLinkSpec: RootLinkSpec{url: "https://www.lianjia.com/city/",
		linkStartFlag: "<ul class=\"city_list_ul\">",
		linkEndFlag:   "<a href=\"https://yw.lianjia.com/\">义乌</a>"},
}

func startMasterEngine() {
	go startMonitor()

	// 房天下
	// rootLinkSpec := RootLinkSpec{url: "https://sh.esf.fang.com/newsecond/esfcities.aspx",
	// 	linkStartFlag: "<div class=\"outCont\" id=\"c02\"",
	// 	linkEndFlag:   "</div>"}

	// 链家
	// rootLinkSpec := RootLinkSpec{url: "https://www.lianjia.com/city/",
	// 	linkStartFlag: "<ul class=\"city_list_ul\">",
	// 	linkEndFlag:   "<a href=\"https://yw.lianjia.com/\">义乌</a>"}

	go startEngineDriver(engineDriver.rootLinkSpec)

	go startEngine()

	startMasterEngineServer()
}

var slavesChan = make(chan string, 10)

func startMasterEngineServer() {
	var port = "8000"
	fmt.Println("监控页面: http://localhost:" + port + "/monitor")

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		slave := r.PostFormValue("s")
		slavesChan <- slave
		registeredSlaves[slave] = struct{}{}
		fmt.Println("加入任务节点:" + slave)
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		var finishTasks []string
		json.Unmarshal([]byte(r.FormValue("fc")), &finishTasks)

		for _, taskName := range finishTasks {
			progress := crawlingTasks[taskName]
			if progress != nil {
				progress.CrawledRecordCount += getPageSize()
				// atomic.AddInt32(*progress.CrawledRecordCount, PageSize)
				incrementPageNum(taskName)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var registeredSlaves = make(map[string]struct{})

func getRegisteredSlaves() []string {
	var slaves = []string{}
	for k, _ := range registeredSlaves {
		slaves = append(slaves, k)
	}
	return slaves
}
