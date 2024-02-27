package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	// "strings"
	"testing"
	// "time"
)

// func readPage() string {
// 	file, _ := os.Open("./page.txt")
// 	defer file.Close()
// 	reader := bufio.NewReader(file)
// 	var page bytes.Buffer
// 	for {
// 		line, _, err := reader.ReadLine()
// 		if err == io.EOF {
// 			break
// 		}
// 		page.Write(line)
// 	}
// 	return page.String()
// }

func TestProtocolAndHost(t *testing.T) {
	link := LinkInfo{"https://esf.fang.com/house/i32/", "dfs"}
	println(link.ProtocolAndHost())
}

func TestJson(t *testing.T) {
	p := CrawlProgress{
		TaskName:           "广州",
		PageNum:            1,
		PageCount:          10,
		HandleThread:       "thread",
		CrawledRecordCount: 98,
		Link:               "http:123",
	}
	fmt.Println(p)
	jsonByte, _ := json.Marshal(p)
	jsonStr := string(jsonByte)
	fmt.Println("-------------" + jsonStr)
}

// func TestPageCount(t *testing.T) {
// file, _ := os.Open("./page.txt")
// defer file.Close()
// reader := bufio.NewReader(file)
// var page bytes.Buffer
// for {
// 	line, _, err := reader.ReadLine()
// 	if err == io.EOF {
// 		break
// 	}
// 	page.Write(line)
// }
// count := getPageCount(page.String())
// fmt.Print("page count --------> ")
// fmt.Println(count)
// }

// func TestDB(t *testing.T) {
// 	count := getMaxIdOfHouse()
// 	fmt.Print("###################count: ")
// 	fmt.Println(count)
// }

// func TestGetPageContentBySearchPage(t *testing.T) {
// u := "http://search.fang.com/captcha-294815364897324240/redirect?h=https://yz.esf.fang.com/house-a01077/"
// url := getPageContentBySearchPage(u)
// fmt.Println("###################url: " + url)
// }

func TestPageSize(t *testing.T) {
	file, _ := os.Open("./page.txt")
	defer file.Close()
	reader := bufio.NewReader(file)
	var page bytes.Buffer
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		page.Write(line)
	}
	s := *getContainKeywordLinks(page.String(), []string{"chushou"})
	fmt.Println(len(s))
}

// func TestCheckRedirectCity(t *testing.T) {
// 	linkInfo := checkRedirectCity(readPage(), LinkInfo{desc: "喧哗", link: ""})
// 	fmt.Println(linkInfo)

// 	a := "/chushou/10_471031220.htm?channel=2,2&psid=1_58_60"
// 	fmt.Println(a[:strings.Index(a, "?")])
// }

// func TestSaveSpeed(t *testing.T) {

// 	house := House{City: "test"}
// 	go func() {
// 		for {
// 			save(house)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			save(house)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			save(house)
// 		}
// 	}()

// 	count := 0
// 	for {
// 		c := getMaxIdOfHouse()
// 		fmt.Printf("count: %d | speed: %d \n", count, c-count)
// 		count = c
// 		time.Sleep(time.Duration(1) * time.Second)
// 	}
// }

func TestSaveSpeed(t *testing.T) {
	// var engineDriver = LianJiaEngineDriver{}
	// for _, link := range engineDriver.firstPageLinks() {
	// 	fmt.Println(link)
	// }

	getPage("https://cz.lianjia.com/ershoufang")
}
