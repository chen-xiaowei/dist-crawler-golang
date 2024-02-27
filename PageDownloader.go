package main

import (
	"io/ioutil"
	// "log"
	"runtime/debug"
	"strings"

	// "bytes"
	"compress/gzip"
	"context"
	"fmt"

	// "io"

	// "os"

	// "runtime/debug"

	// "net"
	"net/http"
	"time"
)

var visitFailedLinks []string

var cookie string

func requestByGet(link string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("user-agent", USER_AGENT)
	req.Header.Add("Cookie", cookie)
	// req.Header.Add("Referer", link)

	// httpclient := &http.Client{
	// 	Transport: &http.Transport{
	// 		Dial: (&net.Dialer{
	// 			Timeout:   30 * time.Second,
	// 			KeepAlive: 30 * time.Second,
	// 		}).Dial,
	// 		TLSHandshakeTimeout:   10 * time.Second,
	// 		ResponseHeaderTimeout: 10 * time.Second,
	// 		ExpectContinueTimeout: 1 * time.Second,
	// 	}}
	ctx, _ := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	req.WithContext(ctx)

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err == nil {
		for _, c := range resp.Header.Values("Set-Cookie") {
			c = c[:strings.Index(c, ";")+1]
			cookie += c + " "
		}
	}
	if resp.StatusCode == 301 || resp.StatusCode == 302 {
		redirectLink := resp.Header.Get("Location")
		// fmt.Println(redirectLink)
		return requestByGet(redirectLink)
	}
	return resp, err
	// return httpclient.Do(req)
}

func getPage(link string) string {
	// RETRY:
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		fmt.Println("http.NewRequest异常|" + err.Error())
		return ""
	}
	req.Header.Add("Accept-Encoding", "gzip")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("user-agent", USER_AGENT)
	req.Header.Add("Cookie", cookie)
	// req.Header.Add("Referer", link)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		debug.PrintStack()
		Err("http.DefaultTransport.RoundTrip异常|" + err.Error() + "|" + link)
		// goto RETRY
		return ""
	}
	for _, c := range resp.Header.Values("Set-Cookie") {
		c = c[:strings.Index(c, ";")+1]
		cookie = c + " "
	}
	if resp.StatusCode == 301 {
		redirectLink := resp.Header.Get("Location")
		return engineDriver.http301Handler()(redirectLink)
	}
	if resp.StatusCode == 302 {
		redirectLink := resp.Header.Get("Location")
		return engineDriver.http302Handler()(redirectLink)
	}
	if resp.StatusCode == 400 {
		fmt.Println("==========cookie=============")
		Err("Bad request for " + link + "\n")
		fmt.Println("cookie: " + cookie)
		fmt.Println("==========cookie=============\n")
		return ""
	}
	defer resp.Body.Close()

	var pageBody []byte
	if "gzip" == resp.Header.Get("Content-Encoding") {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			Err("gzip.NewReader error | " + link)
			return ""
		}
		pageBody, _ = ioutil.ReadAll(reader)
		defer reader.Close()
	} else {
		pageBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			Err("ioutil.ReadAll error | " + link)
			return ""
		}
	}
	return string(pageBody)
}
