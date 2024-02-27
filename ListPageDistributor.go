package main

import (
	"fmt"
	// "strconv"

	// "math/rand"
	"net/http"
	"net/url"
	// "runtime"
	// "sync"
	// "time"
)

var FirstPageLinkChannel = make(chan LinkInfo, 10)

var deliverTaskChannel = make(chan CrawlTask, 10)
var CrawlTasklimiter = make(chan struct{}, 8)

func startCrawlTaskDeliverWork() {
	go func() {
		for crawlTask := range deliverTaskChannel {
			CrawlTasklimiter <- struct{}{}
			if !deliverCrawlTaskToSlave(crawlTask) {
				crawlTaskChannel <- crawlTask
			}
		}
	}()

	for firstPageLinkInfo := range FirstPageLinkChannel {

		go func(firstPageLinkInfo LinkInfo) {
			if firstPageLinkInfo.desc == "" && firstPageLinkInfo.link == "" {
				return
			}
			pageContent := getPage(firstPageLinkInfo.link)
			if "" == pageContent {
				return
			}
			engineDriver.buildCrawlTask(deliverTaskChannel,
				LinkInfo{link: firstPageLinkInfo.link,
					desc:      firstPageLinkInfo.desc,
					pageCount: engineDriver.pageCount(pageContent),
				})
		}(firstPageLinkInfo)
	}
}

func deliverCrawlTaskToSlave(crawlTask CrawlTask) bool {
	if len(slavesChan) == 0 {
		return false
	}
	ipPort, ok := <-slavesChan
	if !ok {
		return false
	}
	deliverTaskUrl := "http://" + ipPort + "/task/accept"

	resp, err := http.PostForm(deliverTaskUrl, url.Values{
		"tn":  {crawlTask.taskName},
		"url": {crawlTask.url},
		"ds":  {crawlTask.desc},
		"pc":  {string(crawlTask.pageCount)},
		"pp":  {crawlTask.pagingPrefix},
	})
	if err != nil {
		fmt.Println("分发爬取任务异常|" + err.Error())
		delete(registeredSlaves, ipPort)
		return false
	}
	slavesChan <- ipPort

	registeredSlaves[ipPort] = struct{}{}

	defer resp.Body.Close()
	if "y" == resp.Header.Get("ac") {
		return true
	}
	return false
}
