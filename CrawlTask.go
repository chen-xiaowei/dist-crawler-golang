package main

import (
	"strings"
)

type CrawlTask struct {
	taskName       string
	desc           string
	url            string
	currentPageNum int
	pageCount      int
	pagingPrefix   string
}

func (crawlTask CrawlTask) ProtocolAndHost() string {
	if crawlTask.url == "" {
		return ""
	}
	idx := strings.Index(crawlTask.url, "//")

	protocol := crawlTask.url[:idx+2]
	host := crawlTask.url[idx+2:]
	idx = strings.Index(host, "/")
	if idx == -1 {
		return crawlTask.url
	}
	host = host[:idx]

	return protocol + host
}
