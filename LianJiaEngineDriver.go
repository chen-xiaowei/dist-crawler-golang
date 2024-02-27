package main

import (
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------//
// --------------------------- 链家 --------------------------------------//
// ----------------------------------------------------------------------//

type LianJiaEngineDriver struct {
	rootLinkSpec RootLinkSpec
}

type LianJiaDetailPageParser struct{}

func (engineDriver LianJiaEngineDriver) nextPageLink() string {
	return "/pg"
}

func (engineDriver LianJiaEngineDriver) rebuildLink(link string) string {
	return link + "/ershoufang"
}

func (engineDriver LianJiaEngineDriver) pageCount(pageContent string) int {
	keyword := "totalPage\":"
	idx := strings.LastIndex(pageContent, keyword)
	if idx == -1 {
		return -1
	}
	strPageCount := pageContent[idx+len(keyword):]
	strPageCount = strPageCount[:strings.Index(strPageCount, ",")]
	pageCount, _ := strconv.Atoi(strPageCount)

	return pageCount
}

func (engineDriver LianJiaEngineDriver) pageSize() int {
	return 30
}

func (engineDriver LianJiaEngineDriver) parseDetailPage(taskName string, detailPageURL string) bool {
	return LianJiaDetailPageParser{}.parse(taskName, detailPageURL)
}

func (engineDriver LianJiaEngineDriver) detailLinkExtractSpecification() DetailLinkExtractSpecification {

	return DetailLinkExtractSpecification{detailLinkStartFlag: "<ul class=\"sellListContent",
		detailLinkEndFlag: "<div id=\"noResultPush\"",
		keywordsInLink:    []string{"ershoufang", "html"}}
}

func (engineDriver LianJiaEngineDriver) buildCrawlTask(deliverTaskChannel chan CrawlTask, firstPageLinkInfo LinkInfo) {
	pageContent := getPage(firstPageLinkInfo.link)

	regionStart := "<div data-role=\"ershoufang\""
	idx := strings.Index(pageContent, regionStart)
	if idx == -1 {
		return
	}
	pageCount := engineDriver.pageCount(pageContent)
	if pageCount == -1 {
		return
	}
	if firstPageLinkInfo.pageCount < 100 {
		deliverTaskChannel <- CrawlTask{taskName: firstPageLinkInfo.desc,
			url:          firstPageLinkInfo.link,
			pageCount:    pageCount,
			pagingPrefix: "/pg"}
		return
	}
	regionLinksRange := pageContent[idx+len(regionStart):]
	regionLinksRange = regionLinksRange[:strings.Index(regionLinksRange, "</div>")]

	for _, linkInfo := range extractLinks(regionLinksRange) {

		pageContent = getPage(firstPageLinkInfo.ProtocolAndHost() + linkInfo.link)
		idx = strings.Index(pageContent, regionStart)
		if idx == -1 {
			continue
		}
		subRegionLinksRange := pageContent[idx+len(regionStart):]
		subRegionLinksRange = subRegionLinksRange[strings.Index(subRegionLinksRange, "</div>")+6:]
		subRegionLinksRange = subRegionLinksRange[:strings.Index(subRegionLinksRange, "</div>")+6]

		// fmt.Println("===========================================================")
		// fmt.Println(subRegionLinksRange)
		// fmt.Println("===========================================================")

		for _, subLinkInfo := range extractLinks(subRegionLinksRange) {
			url := firstPageLinkInfo.ProtocolAndHost() + subLinkInfo.link

			pageCount := engineDriver.pageCount(getPage(url))

			// fmt.Printf("--------%s %d \n", url, pageCount)

			if pageCount == -1 {
				continue
			}
			crawlTask := CrawlTask{taskName: firstPageLinkInfo.desc + "-" + linkInfo.desc + "-" + subLinkInfo.desc,
				url:          url,
				pageCount:    pageCount,
				desc:         linkInfo.desc + "-" + subLinkInfo.desc,
				pagingPrefix: "/pg"}

			deliverTaskChannel <- crawlTask
		}
	}
}

func (engineDriver LianJiaEngineDriver) http301Handler() func(url string) string {
	return func(url string) string {
		return getPage(url)
	}
}

func (engineDriver LianJiaEngineDriver) http302Handler() func(url string) string {
	return func(url string) string {
		return ""
	}
}
