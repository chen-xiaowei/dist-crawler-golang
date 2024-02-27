package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// import . "crawler-framework/main"

// ------------------------------------------------------------------------//
// --------------------------- 房天下 --------------------------------------//
// ------------------------------------------------------------------------//

type FangTianXiaEngineDriver struct {
	rootLinkSpec RootLinkSpec
}

type FangTianXiaDetailPageParser struct{}

type FangTianXiaDetailPageLinkExtractor struct{}

func (engineDriver FangTianXiaEngineDriver) buildCrawlTask(deliverTaskChannel chan CrawlTask, firstPageLinkInfo LinkInfo) {

	// fmt.Printf("%#v \n", firstPageLinkInfo)

	pageContent := getPage(firstPageLinkInfo.link)
	if noHouseSource(pageContent, firstPageLinkInfo) {
		return
	}
	linkInfo, err := checkRedirectCity(pageContent, firstPageLinkInfo)
	if err != nil {
		return
	}
	firstPageLinkInfo = linkInfo

	pageCount := engineDriver.pageCount(pageContent)
	if pageCount == -1 {
		return
	}
	// fmt.Printf("pageCount - %d \n", firstPageLinkInfo.pageCount)

	if firstPageLinkInfo.pageCount < 100 {
		deliverTaskChannel <- CrawlTask{taskName: firstPageLinkInfo.desc,
			url:          firstPageLinkInfo.link,
			pageCount:    pageCount,
			pagingPrefix: "/house/i3"}
		return
	}
	//如果 页数 等于100 则根据 子区域 为单位构建 crawl task
	regionStart := "<ul class=\"clearfix choose_screen floatl\">"
	idx := strings.Index(pageContent, regionStart)
	if idx == -1 {
		return
	}
	regionLinksRange := pageContent[idx+len(regionStart):]
	regionEnd := "</ul>"
	regionLinksRange = regionLinksRange[:strings.Index(regionLinksRange, regionEnd)]

	subRegionStart := "<li class=\"area_sq\">"
	subRegionEnd := "</ul>"

	for _, linkInfo := range extractLinks(regionLinksRange) {
		pageContent = getPage(firstPageLinkInfo.ProtocolAndHost() + linkInfo.link)

		idx := strings.Index(pageContent, subRegionStart)
		if idx == -1 {
			continue
		}
		subRegionLinksRange := pageContent[idx+len(subRegionStart):]
		subRegionLinksRange = subRegionLinksRange[:strings.Index(subRegionLinksRange, subRegionEnd)+len(subRegionEnd)]

		for _, subLinkInfo := range extractLinks(subRegionLinksRange) {
			pageCount := engineDriver.pageCount(getPage(firstPageLinkInfo.ProtocolAndHost() + subLinkInfo.link))
			if pageCount == -1 {
				continue
			}
			taskName := firstPageLinkInfo.desc + "-" + linkInfo.desc + "-" + subLinkInfo.desc
			deliverTaskChannel <- CrawlTask{taskName: taskName,
				url:          firstPageLinkInfo.ProtocolAndHost() + subLinkInfo.link,
				pageCount:    pageCount,
				desc:         linkInfo.desc + "-" + subLinkInfo.desc,
				pagingPrefix: "/i3"}
		}
	}
}

func getPageSize() int {
	return engineDriver.pageSize()
}

func (engineDriver FangTianXiaEngineDriver) pageSize() int {
	return 60
}

func (engineDriver FangTianXiaEngineDriver) pageCount(pageContent string) int {
	idx := strings.Index(pageContent, ">共")
	if idx == -1 {
		return 1
	}
	pageCount := pageContent[idx+4:]
	endIdx := strings.Index(pageCount, "页<")
	if endIdx == -1 {
		pageCount = pageCount[:1]
	} else {
		pageCount = pageCount[:endIdx]
	}
	count, _ := strconv.Atoi(pageCount)

	return count
}

func (engineDriver FangTianXiaEngineDriver) http302Handler() func(url string) string {
	return func(url string) string {
		if !strings.Contains(url, "search") {
			return getPage(url)
		}
		searchPage := getPage(url)
		if searchPage == "" {
			Err("未获取到[搜索跳转页]|" + url)
			return ""
		}
		//页面图灵人机测试
		authPageContent := "请拖动滑块进行验证"
		if strings.Contains(searchPage, authPageContent) {
			sleepSec := 300
			alert := "******************************************\n"
			alert += "********访问地址要求图灵人机测试***************\n"
			alert += "url:" + url + "\n"
			alert += "将在" + strconv.Itoa(sleepSec) + "秒后进行重试\n"
			alert += "********" + time.Now().Format("2007/01/02 03:04:05 PM Mon") + "*******\n"
			alert += "******************************************\n"

			fmt.Println(alert)

			time.Sleep(time.Duration(sleepSec) * time.Second)

			getPage(url)
		}
		//这是搜索页为防爬虫在JS中定义的url变量，将t3 t4拼接后再次访问即为真正详情页
		detailPageUrlPosition := "var t4='"
		paramPosition := "var t3='"
		t4Idx := strings.Index(searchPage, detailPageUrlPosition)
		detailUrl := searchPage[t4Idx+len(detailPageUrlPosition):]
		detailUrl = detailUrl[:strings.Index(detailUrl, "'")]

		param := searchPage[strings.Index(searchPage, paramPosition)+len(paramPosition):]
		param = param[:strings.Index(param, "'")]

		return getPage(detailUrl + "?" + param)
	}
}

func (engineDriver FangTianXiaEngineDriver) rebuildLink(link string) string {
	return link
}

func (engineDriver FangTianXiaEngineDriver) http301Handler() func(url string) string {
	return func(url string) string {
		return getPage(url)
	}
}

func (engineDriver FangTianXiaEngineDriver) parseDetailPage(taskName string, detailPageURL string) bool {
	return FangTianXiaDetailPageParser{}.parse(taskName, detailPageURL)
}

func (engineDriver FangTianXiaEngineDriver) detailLinkExtractSpecification() DetailLinkExtractSpecification {

	return DetailLinkExtractSpecification{detailLinkStartFlag: "<div class=\"shop_list shop_list_4\">",
		detailLinkEndFlag: ">共",
		keywordsInLink:    []string{"chushou"}}
}

func checkRedirectCity(pageContent string, linkInfo LinkInfo) (LinkInfo, error) {
	//如果城市首页经过跳转, 则用跳转后的城市和对应链接代替被跳转城市
	ordinaryCity := linkInfo.desc
	cityNameStartFeature := "<a href=\"#\">"
	idx := strings.Index(pageContent, cityNameStartFeature)
	if idx == -1 {
		Err("未在[" + linkInfo.desc + "]首页获取到跳转城市" + "|" + linkInfo.link)
		return linkInfo, errors.New("can not find city name flag")
	}
	redirectCity := pageContent[idx+len(cityNameStartFeature):]
	redirectCity = redirectCity[:strings.Index(redirectCity, "<")]
	if redirectCity == ordinaryCity {
		return linkInfo, nil
	}
	redirectCity = redirectCity + "-" + ordinaryCity
	fmt.Println("城市跳转 " + ordinaryCity + " -> " + redirectCity + "[" + linkInfo.link + "]")

	return LinkInfo{link: linkInfo.link, desc: redirectCity}, nil
}

func noHouseSource(listPageContent string, linkInfo LinkInfo) bool {
	sorry := strings.Index(listPageContent, "很抱歉，")
	notFound := strings.Index(listPageContent, "没有找到")
	houseSource := strings.Index(listPageContent, "相符的房源！")
	if sorry != -1 && notFound != -1 && houseSource != -1 {
		if houseSource > notFound && notFound > sorry {
			fmt.Println("[" + linkInfo.desc + "]无房源|" + linkInfo.link)
			return true
		}
	}
	return false
}
