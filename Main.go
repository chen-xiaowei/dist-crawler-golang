package main

import (
	"flag"
	"fmt"
)

var start string
var Master string

func init() {
	flag.StringVar(&start, "start", "m", "引擎启动模式 m:master / s:slave")
	flag.StringVar(&Master, "master", "", "Master's ip and port(ip:port)")
}

func main() {
	flag.Parse()
	if start == "m" {
		startMasterEngine()
	} else if start == "s" {
		if Master == "" {
			fmt.Println("Please input Master's ip and port(ip:port)")
			return
		}
		startSlaveEngine()
	} else {
		fmt.Println("Please input param for -start with m or s")
	}
}

func isMasterMode() bool {
	return start == "m"
}
