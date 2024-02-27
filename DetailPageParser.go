package main

import (
	"fmt"
	// "fmt"
	"strings"
	// "sync"
)

type DetailPageParser interface {
	parse(taskName string, url string) bool
}

func (parser FangTianXiaDetailPageParser) parse(taskName string, url string) bool {
	pageContent := getPage(url)
	idx := strings.Index(pageContent, "lab\">区")
	if idx == -1 {
		// 提示页面: 抱歉，未找到相关页面
		Err("未在当前页面找到指定值|" + url)
		return false
	}
	remainContent := pageContent[idx:]
	idx = strings.Index(remainContent, "</a>")
	region := remainContent[:idx]
	i := len(region) - 1
	for ; i > 0; i-- {
		if region[i] == 62 { //62 -> ">"
			break
		}
	}
	var house House
	house.Region = strings.TrimSpace(region[i+1:])

	subRegion := remainContent[idx+4:]
	subRegion = subRegion[:strings.Index(subRegion, "</a>")]
	i = len(subRegion) - 1
	for ; i > 0; i-- {
		if subRegion[i] == 62 { //62 -> ">"
			break
		}
	}
	house.SubRegion = strings.TrimSpace(subRegion[i+1:])

	house.BuildYear = strings.TrimRight(getFieldValue(pageContent, ">建筑年代<"), "年")
	house.BuildingCount = getFieldValue(pageContent, "总楼栋数</span>")
	house.HouseholdCount = getFieldValue(pageContent, "总  户  数</span>")
	house.PublishTime = getFieldValue(pageContent, "挂牌时间</span>")

	idx = strings.Index(pageContent, ">建筑面积<")
	if idx != -1 {
		acreage := pageContent[:idx]
		acreage = acreage[:strings.LastIndex(acreage, "平米")]
		acreage = acreage[strings.LastIndex(acreage, ">")+1:]
		house.Acreage = acreage
	}

	if !strings.Contains(pageContent, "面议") && !strings.Contains(pageContent, "--") {
		unitPrice := pageContent[:strings.Index(pageContent, "单价</div>")]
		idx = strings.LastIndex(unitPrice, "元")
		if idx != -1 {
			unitPrice = unitPrice[:idx]
			unitPrice = unitPrice[strings.LastIndex(unitPrice, ">")+1:]
			house.UnitPrice = unitPrice
		}
	}
	idx = strings.Index(pageContent, "小&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;区")
	if idx == -1 {
		idx = strings.Index(pageContent, "小 区")
	}
	if idx != -1 {
		address := pageContent[idx:]
		address = address[:strings.Index(address, "</a")]
		address = address[strings.LastIndex(address, ">")+1:]
		house.Address = strings.TrimSpace(address)

		if house.Address == "l>" {
			fmt.Println(url)
		}
	}
	idx = strings.Index(taskName, "-")
	if idx != -1 {
		house.City = taskName[:idx]
	} else {
		house.City = taskName
	}
	house.Link = url

	save(house)

	return true
}

func getFieldValue(pageContent string, key string) string {
	idx := strings.Index(pageContent, key)
	if idx == -1 {
		return ""
	}
	value := pageContent[idx+len(key):]
	value = value[:strings.Index(value, "</")]

	i := len(value) - 1
	//从后向前查找 ">"
	for ; i > 0; i-- {
		if value[i] == 62 { //62 -> ">"
			break
		}
	}
	value = strings.TrimSpace(value[i+1:])
	return value
}

// ------------------------------ 链家------------------------------//

func (parser LianJiaDetailPageParser) parse(taskName string, detailPageURL string) bool {
	detailPageContent := getPage(detailPageURL)
	if detailPageContent == "" {
		return false
	}
	var house House

	keyword := "建筑面积</span>"
	idx := strings.Index(detailPageContent, keyword)
	if idx == -1 {
		Err("未定位到keyword:[" + keyword + "]|" + detailPageURL)
		return true
	}
	house.Acreage = detailPageContent[idx+len(keyword):]
	idx = strings.Index(house.Acreage, "㎡")
	house.Acreage = house.Acreage[:idx]

	keyword = "挂牌时间</span>"
	sIdx := strings.Index(detailPageContent, keyword)
	house.PublishTime = strings.Trim(detailPageContent[sIdx+len(keyword):], "\n")
	house.PublishTime = strings.Trim(house.PublishTime, " ")
	house.PublishTime = house.PublishTime[6:strings.Index(house.PublishTime, "</")]

	keyword = "class=\"unitPriceValue\">"
	sIdx = strings.Index(detailPageContent, keyword)
	house.UnitPrice = detailPageContent[sIdx+len(keyword):]
	house.UnitPrice = house.UnitPrice[:strings.Index(house.UnitPrice, "<")]

	keyword = "class=\"info \">"
	sIdx = strings.Index(detailPageContent, keyword)
	house.Address = detailPageContent[sIdx+len(keyword):]
	house.Address = house.Address[0:strings.Index(house.Address, "<")]

	keyword = "target=\"_blank\">"
	sIdx = strings.Index(detailPageContent, keyword)
	house.Region = detailPageContent[sIdx+len(keyword):]
	sIdx = strings.Index(house.Region, keyword)
	house.Region = house.Region[sIdx+len(keyword):]
	sIdx = strings.Index(house.Region, keyword)
	house.Region = house.Region[sIdx+len(keyword):]

	sIdx = strings.Index(house.Region, keyword)
	house.SubRegion = house.Region[sIdx+len(keyword):]
	house.Region = house.Region[:strings.Index(house.Region, "<")]
	house.SubRegion = house.SubRegion[0:strings.Index(house.SubRegion, "<")]

	idx = strings.Index(taskName, "-")
	if idx != -1 {
		house.City = taskName[:idx]
	} else {
		house.City = taskName
	}
	house.Link = detailPageURL

	save(house)

	return true
}
