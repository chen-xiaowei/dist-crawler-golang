package main

import (
	"strings"
)

type LinkInfo struct {
	link      string
	desc      string
	pageCount int
}

// http(s)://host
func (linkInfo LinkInfo) ProtocolAndHost() string {
	if linkInfo.link == "" {
		return ""
	}
	idx := strings.Index(linkInfo.link, "//")

	protocol := linkInfo.link[:idx+2]
	host := linkInfo.link[idx+2:]
	idx = strings.Index(host, "/")
	if idx == -1 {
		return linkInfo.link
	}
	host = host[:idx]

	return protocol + host
}
