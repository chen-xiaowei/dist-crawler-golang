package main

import (
	// "fmt"
	"strings"
	// "sync"
)

type DetailLinkExtractSpecification struct {
	detailLinkStartFlag string
	detailLinkEndFlag   string
	keywordsInLink      []string
}

func extractDetailPageLink(spec DetailLinkExtractSpecification, crawlTask CrawlTask) {
	// var detailPageLimiter = make(chan struct{}, 5)

	for {
		nextPageURL := nextPageLink(&crawlTask)
		if nextPageURL == "" {
			break
		}
		// go func() {
		incrementPageNum(crawlTask.taskName)

		for _, detailPageLink := range *getLinksInPage(nextPageURL, spec) {
			// detailPageLimiter <- struct{}{}

			if strings.Contains(detailPageLink, "http") {
				engineDriver.parseDetailPage(crawlTask.taskName, detailPageLink)
			} else {
				engineDriver.parseDetailPage(crawlTask.taskName, crawlTask.ProtocolAndHost()+detailPageLink)
			}
			incrementCrawledRecordCount(crawlTask)

			// <-detailPageLimiter
		}
		// }()
	}

	<-CrawlTasklimiter

	// // 从列表页抽取详情页链接
	// 如果当前节点以Master模式启动，则自增本地pageNum
	// 否则说明当前节点是Slave，则将完成的城市发送给Master进行爬取进度同步
	// if isMasterMode() {
	// 	incrementPageNum(crawlTask.taskName)
	// } else {
	// 	finishedTask = append(finishedTask, crawlTask.taskName)
	// }

}

func getLinksInPage(url string, spec DetailLinkExtractSpecification) *[]string {
	var links []string

	listPageContent := getPage(url)
	idx := strings.Index(listPageContent, spec.detailLinkStartFlag)
	if idx == -1 {
		return &links
	}
	listPageContent = listPageContent[idx:]
	idx = strings.Index(listPageContent, spec.detailLinkEndFlag)
	if idx == -1 {
		return &links
	}
	//包含详情页URI的列表页的部分
	listPageContent = listPageContent[:idx]

	return getContainKeywordLinks(listPageContent, spec.keywordsInLink)
}

func getContainKeywordLinks(listPageContent string, keywords []string) *[]string {
	var links []string

	lastURI := ""
	var idx int
	for strings.Contains(listPageContent, "href=") {
		idx = strings.Index(listPageContent, "<a ")
		listPageContent = listPageContent[idx:]
		idx = strings.Index(listPageContent, ">")
		if idx == -1 {
			continue
		}
		aTag := listPageContent[:idx+1]

		hrefIdx := strings.Index(aTag, "href=\"")
		if hrefIdx == -1 {
			listPageContent = listPageContent[idx:]
			continue
		}
		href := aTag[hrefIdx+6:]

		if strings.Index(href, "\"") == -1 {
			listPageContent = listPageContent[len(href):]
			continue
		}
		href = href[:strings.Index(href, "\"")]

		containKeyword := true
		for _, keyword := range keywords {
			if !strings.Contains(aTag, keyword) {
				containKeyword = false
			}
		}
		if containKeyword && lastURI != href {
			questionIdx := strings.Index(href, "?")
			if questionIdx != -1 {
				href = href[:questionIdx]
			}
			lastURI = href

			var param string
			chanIdx := strings.Index(aTag, "channel=")
			if chanIdx != -1 {
				attrChannel := aTag[chanIdx+8:]
				attrChannel = attrChannel[:strings.Index(attrChannel, " ")]

				param += "channel=" + attrChannel + "&"
			}

			psIdx := strings.Index(aTag, "ps=\"")
			if psIdx != -1 {
				attrPS := aTag[psIdx+4:]
				attrPS = attrPS[:strings.Index(attrPS, "\"")]

				param += "psid=" + attrPS
			}
			if param != "" {
				href += "?" + param
			}
			links = append(links, href)
		}
		listPageContent = listPageContent[idx:]
	}
	return &links
}
