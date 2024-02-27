package main

import (
	// "time"
	// "io/ioutil"
	// "net/http"
	"fmt"
	"strings"
)

// func init() {
// 	go startMonitor()
// }

var PageSize int

func getPageSize() int {
	return PageSize
}

var ConcurrentCityLimiter = make(chan struct{}, CONCURRENT_CITY_NUM)

func startEngineDriver() {
	fmt.Println("Engine Driver started")

	rootURL := "https://sh.esf.fang.com/newsecond/esfcities.aspx"

	pageContent := getPage(rootURL)
	if pageContent == "" {
		return
	}
	htmlTagOfLinkIn := "<div class=\"outCont\" id=\"c02\""
	linkStartIdx := strings.Index(pageContent, htmlTagOfLinkIn)
	pageContent = pageContent[linkStartIdx:]
	linkEndIdx := strings.Index(pageContent, "</div>")
	pageContent = pageContent[:linkEndIdx]

	go startListPageDeliverWork()

	// 抽取城市首页Link
	for strings.Contains(pageContent, "href=") {

		idx := strings.Index(pageContent, "href=\"")
		pageContent = pageContent[idx+6:]
		idx = strings.Index(pageContent, "\"")
		cityFirstPageLink := pageContent[:idx]
		pageContent = pageContent[idx:]
		pageContent = pageContent[strings.Index(pageContent, ">")+1:]
		city := pageContent[:strings.Index(pageContent, "<")]

		if strings.HasSuffix(cityFirstPageLink, "/") {
			cityFirstPageLink = cityFirstPageLink[:len(cityFirstPageLink)-1]
		}
		if strings.HasPrefix(cityFirstPageLink, "//") {
			cityFirstPageLink = "https:" + cityFirstPageLink
		}
		if getPageSize() == 0 {
			page := getPage(cityFirstPageLink)
			if strings.Contains(page, "下一页") {
				PageSize = len(*getContainKeywordLinks(page, "chushou"))
			}
		}

		linkInfo := LinkInfo{cityFirstPageLink, city}
		// linkInfo.desc = "张北"
		// linkInfo.link = "https://zhangbei.esf.fang.com"

		ConcurrentCityLimiter <- struct{}{}

		CityFirstPageLinkChannel <- linkInfo
	}
	CityFirstPageLinkChannel <- LinkInfo{"", ""}

	var input string
	fmt.Scanln(&input)
}
