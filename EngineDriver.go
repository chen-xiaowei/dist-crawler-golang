package main

import (
	"strconv"
	// "errors"
	// "fmt"
	// "strconv"
	"strings"
	// "time"
)

type EngineDriver interface {
	buildCrawlTask(crawlTaskChannel chan CrawlTask, firstPageLink LinkInfo)
	pageSize() int
	pageCount() int
	parseDetailPage(taskName string, url string) bool
	nextPageLink() string
	http301Handler() func() string
	http302Handler() func() string
}

type RootLinkSpec struct {
	url           string
	linkStartFlag string
	linkEndFlag   string
}

func startEngineDriver(rootLinkSpec RootLinkSpec) {
	go startCrawlTaskDeliverWork()

	extractLinksByRootLinkSpec(rootLinkSpec)
	// var input string
	// fmt.Scanln(&input)
}

func extractLinksByRootLinkSpec(rootLinkSpec RootLinkSpec) {
	pageContent := getPage(rootLinkSpec.url)
	if pageContent == "" {
		return
	}
	linkStartIdx := strings.Index(pageContent, rootLinkSpec.linkStartFlag)
	pageContent = pageContent[linkStartIdx:]
	linkEndIdx := strings.Index(pageContent, rootLinkSpec.linkEndFlag)
	pageContent = pageContent[:linkEndIdx]

	for strings.Contains(pageContent, "href=") {

		idx := strings.Index(pageContent, "href=\"")
		pageContent = pageContent[idx+6:]
		idx = strings.Index(pageContent, "\"")
		firstPageLinkOfSubject := pageContent[:idx]
		pageContent = pageContent[idx:]
		pageContent = pageContent[strings.Index(pageContent, ">")+1:]
		hrefText := pageContent[:strings.Index(pageContent, "<")]

		if strings.HasSuffix(firstPageLinkOfSubject, "/") {
			firstPageLinkOfSubject = firstPageLinkOfSubject[:len(firstPageLinkOfSubject)-1]
		}
		if strings.HasPrefix(firstPageLinkOfSubject, "//") {
			firstPageLinkOfSubject = "https:" + firstPageLinkOfSubject
		}
		resp, err := requestByGet(firstPageLinkOfSubject)
		if err == nil {
			url := resp.Request.URL.String()
			if strings.HasSuffix(url, "/") {
				url = url[:len(url)-1]
			}
			url = engineDriver.rebuildLink(url)
			FirstPageLinkChannel <- LinkInfo{link: url, desc: hrefText}
		}
	}
}

// 公用方法
func extractLinks(pageContent string) []LinkInfo {
	var links []LinkInfo

	for strings.Contains(pageContent, "href=") {
		idx := strings.Index(pageContent, "href=\"")
		pageContent = pageContent[idx+6:]
		idx = strings.Index(pageContent, "\"")
		firstPageLink := pageContent[:idx]
		pageContent = pageContent[idx:]
		pageContent = pageContent[strings.Index(pageContent, ">")+1:]
		hrefText := pageContent[:strings.Index(pageContent, "<")]

		if strings.HasSuffix(firstPageLink, "/") {
			firstPageLink = firstPageLink[:len(firstPageLink)-1]
		}
		if strings.HasPrefix(firstPageLink, "//") {
			firstPageLink = "https:" + firstPageLink
		}
		links = append(links, LinkInfo{link: firstPageLink, desc: hrefText})
	}
	return links
}

func nextPageLink(crawlTask *CrawlTask) string {
	if crawlTask.currentPageNum > crawlTask.pageCount {
		return ""
	}
	if crawlTask.currentPageNum == 0 {
		crawlTask.currentPageNum += 2
		return crawlTask.url
	}
	nextURL := crawlTask.url + crawlTask.pagingPrefix + strconv.Itoa(crawlTask.currentPageNum)
	crawlTask.currentPageNum++

	return nextURL
}
