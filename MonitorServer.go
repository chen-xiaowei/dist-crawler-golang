package main

import (
	"bufio"
	"bytes"

	// "encoding/binary"
	"encoding/json"

	// "fmt"
	// "math/rand"
	// "strconv"
	// "time"
	"io"
	"log"
	"net/http"
	"os"
)

func startMonitor() {
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("./monitor.html")
		defer file.Close()
		reader := bufio.NewReader(file)
		var page bytes.Buffer
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			page.Write(line)
			page.WriteString("\n")
		}
		w.Write(page.Bytes())
	})

	http.HandleFunc("/monitor/data", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{}, 5)
		data["crawlingTasks"] = getCrawlingTasks()
		data["finishedCrawlTasks"] = getFinishedTasks()
		data["pageSize"] = getPageSize()
		data["slaves"] = getRegisteredSlaves()
		json, _ := json.Marshal(data)

		// fmt.Println(string(json))

		w.Write(json)
	})

	http.HandleFunc("/monitor/data/crawledCount", func(w http.ResponseWriter, r *http.Request) {
		// println("getCrawledCount " + getCrawledCount())
		w.Write([]byte(getCrawledCount()))

		// data := int64(getCrawledCount())
		// bytebuf := bytes.NewBuffer([]byte{})
		// binary.Write(bytebuf, binary.BigEndian, data)
		// w.Write(bytebuf.Bytes())
	})
	// rand.Seed(time.Now().UnixNano())
	// port := strconv.Itoa(8800 + rand.Intn(30))
	// fmt.Println("#####监控页面端口:[" + port + "]######")
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":8800", nil))
}

type CrawlProgress struct {
	TaskName           string `json:"taskName"`
	PageNum            int    `json:"pageNum"`
	PageCount          int    `json:"pageCount"`
	RecordCount        int    `json:"recordCount"`
	CrawledRecordCount int    `json:"crawledRecordCount"`
	Link               string `json:"link"`
}
