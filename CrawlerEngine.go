package main

// "bytes"
// "container/list"
// "encoding/json"
import (
	"fmt"
	"strconv"
	// "sync/atomic"
)

var crawlTaskChannel = make(chan CrawlTask, 10)

var crawlingTasks = make(map[string]*CrawlProgress, 500)

var finishedTasks []string

func startEngine() {
	fmt.Println("Crawler Engine started")

	for crawlTask := range crawlTaskChannel {
		lastPageURL := crawlTask.url + crawlTask.pagingPrefix + strconv.Itoa(crawlTask.pageCount)
		spec := engineDriver.detailLinkExtractSpecification()
		recordCount := (crawlTask.pageCount-1)*engineDriver.pageSize() + len(*getLinksInPage(lastPageURL, spec))

		addCrawlingTask(
			CrawlProgress{
				TaskName:    crawlTask.taskName,
				Link:        crawlTask.url,
				PageNum:     0,
				PageCount:   crawlTask.pageCount,
				RecordCount: recordCount,
			},
		)
		go extractDetailPageLink(engineDriver.detailLinkExtractSpecification(), crawlTask)
	}
}

func addCrawlTask(crawlTask CrawlTask) bool {
	if len(crawlTaskChannel) > 1 {
		return false
	}
	select {
	case crawlTaskChannel <- crawlTask:
		fmt.Println("接收Master爬取任务[" + crawlTask.taskName + "]->" + crawlTask.url)

		addCrawlingTask(CrawlProgress{
			TaskName:  crawlTask.taskName,
			PageNum:   0,
			PageCount: crawlTask.pageCount,
			// CrawledRecordCount: crawlTask.recordCount,
			Link: crawlTask.url,
		})
		return true
	default:
		return false
	}
}

func getCrawlingTasks() []*CrawlProgress {
	var progress []*CrawlProgress
	for _, v := range crawlingTasks {
		progress = append(progress, v)
	}
	return progress
}

func getFinishedTasks() []string {
	return finishedTasks
}

func getCrawledCount() string {
	return getMaxIdOfHouse()
}

func addCrawlingTask(progress CrawlProgress) {
	crawlingTasks[progress.TaskName] = &progress
}

func finishedCrawlTask(taskName string) {
	delete(crawlingTasks, taskName)
	finishedTasks = append(finishedTasks, taskName)
}

func incrementPageNum(taskName string) {
	progress := crawlingTasks[taskName]
	// fmt.Println(progress)
	if progress != nil && progress.PageNum < progress.PageCount {
		progress.PageNum++
		// atomic.AddInt32(&progress.PageNum, 1)
	}
}

func incrementCrawledRecordCount(crawlTask CrawlTask) {
	progress := crawlingTasks[crawlTask.taskName]
	if progress != nil {
		if progress.CrawledRecordCount == progress.RecordCount-1 {
			finishedCrawlTask(crawlTask.taskName)
			return
		}
		// fmt.Printf("--------- %s %d - %d \n", crawlTask.taskName, progress.CrawledRecordCount, progress.RecordCount)
		progress.CrawledRecordCount++
		// atomic.AddInt32(&progress.CrawledRecordCount, 1)
	}
}
