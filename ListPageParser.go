package main

import (
	"fmt"
	"strings"
	"sync"
)

var listPageChannel = make(chan LinkInfo, 10)

func addCrawlTask(listPageLinkInfo LinkInfo) bool {
	select {
	case listPageChannel <- listPageLinkInfo:
		fmt.Println("接收Master爬取任务[" + listPageLinkInfo.desc + "]->" + listPageLinkInfo.link)

		addCrawlingCity(CrawlProgress{
			City:               listPageLinkInfo.desc,
			PageNum:            0,
			PageCount:          1,
			HandleThread:       "",
			CrawledRecordCount: 0,
			Link:               listPageLinkInfo.link,
		})

		return true
	default:
		return false
	}
}

var listPageParserLimiter = make(chan struct{}, CONCURRENT_LIST_PARSER_NUM)

func startParseListPageWork() {
	for listPageLinkInfo := range listPageChannel {
		//if 表示一个城市的链接已抽取完毕
		if listPageLinkInfo.desc != "" && listPageLinkInfo.link == "" {
			finishedCrawlCity(listPageLinkInfo.desc)
			<-ConcurrentCityLimiter
			continue
		}
		listPageParserLimiter <- struct{}{}

		go extractDetailPageLink(listPageLinkInfo)
	}
}

// 从列表页抽取详情页链接
func extractDetailPageLink(listPageLinkInfo LinkInfo) {

	// 如果当前节点以Master模式启动，则自增本地pageNum
	// 否则说明当前节点是Slave，则将完成的城市发送给Master进行爬取进度同步
	if isMasterMode() {
		incrementPageNum(listPageLinkInfo.desc)
	} else {
		finishedListPageCities = append(finishedListPageCities, listPageLinkInfo.desc)
	}

	listPageContent := getPage(listPageLinkInfo.link)

	//列表开始标签
	listStartFlag := "<div class=\"shop_list shop_list_4\">"
	idx := strings.Index(listPageContent, listStartFlag)
	if idx == -1 {
		return
	}
	listPageContent = listPageContent[idx:]
	listEndFlag := ">共"

	idx = strings.Index(listPageContent, listEndFlag)
	listPageContent = listPageContent[:idx]

	detailPageLinks := *getContainKeywordLinks(listPageContent, "chushou")

	grCtrl := make(chan struct{}, CONCURRENT_DETAIL_PAGE_HANDLER_NUM)

	var wg sync.WaitGroup
	for _, uri := range detailPageLinks {
		grCtrl <- struct{}{}

		// fmt.Println(len(grCtrl))

		wg.Add(1)
		go func(uri string, listPage LinkInfo) {
			saveMessageInDetailPage(listPage.desc, listPage.ProtocolAndHost()+uri)

			<-grCtrl
			wg.Done()
		}(uri, listPageLinkInfo)
	}
	// wg.Wait()

	<-listPageParserLimiter
}

func getContainKeywordLinks(listPageContent string, keyword string) *[]string {
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

		href = href[:strings.Index(href, "\"")]

		if strings.Contains(aTag, keyword) && lastURI != href {
			questionIdx := strings.Index(href, "?")
			if questionIdx != -1 {
				href = href[:questionIdx]
			}
			lastURI = href

			// fmt.Println(aTag)

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
			// fmt.Println(href)

			links = append(links, href)
		}
		listPageContent = listPageContent[idx:]
	}
	return &links
}
