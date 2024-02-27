package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// crawler-framework.exe -start s -master 192.168.1.110:8000

func startSlaveEngine() {
	slavePortChan := make(chan string)

	go startAcceptCrawlTaskServer(slavePortChan)

	go startEngine()

	go register(slavePortChan)

	reportCrawlCircumstanceTask()
}

func register(slavePortChan chan string) {
	taskServerPort := <-slavePortChan
	slaveIPport := localIP() + taskServerPort
	masterRegisterUrl := "http://" + Master + "/register"
	for {
		_, err := http.PostForm(masterRegisterUrl, url.Values{"s": {slaveIPport}})
		if err != nil {
			fmt.Println("向Master注册异常|" + err.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		// defer resp.Body.Close()
		break
	}
	fmt.Println("Slave注册成功")
}

func startAcceptCrawlTaskServer(slavePortChan chan string) {
	http.HandleFunc("/task/accept", func(w http.ResponseWriter, r *http.Request) {
		// link := LinkInfo{link: r.FormValue("url"), desc: r.FormValue("city")}
		pageCount, _ := strconv.Atoi(r.FormValue("pc"))

		task := CrawlTask{
			taskName:     r.FormValue("tn"),
			url:          r.FormValue("url"),
			desc:         r.FormValue("ds"),
			pageCount:    pageCount,
			pagingPrefix: r.FormValue("pp"),
		}
		if addCrawlTask(task) {
			w.Header().Set("ac", "y")
		} else {
			w.Header().Set("ac", "n")
		}
	})
	rand.Seed(time.Now().UnixNano())
	slaveServerPort := ":" + strconv.Itoa(9001+rand.Intn(30))
	slavePortChan <- slaveServerPort

	fmt.Println("Slave Port:" + slaveServerPort)

	log.Fatal(http.ListenAndServe(slaveServerPort, nil))
}

var finishedTask []string

// 将 finishedTask 数组批量发送给master进行爬取进度同步，同步成功后清除该数组
func reportCrawlCircumstanceTask() {
	// ticker := time.NewTicker(3 * time.Second)
	// defer ticker.Stop()

	// for range ticker.C {

	// }

	for {
		if (len(finishedTask)) == 0 {
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		cities, _ := json.Marshal(finishedTask)

		fmt.Print("向Master同步爬取进度")
		fmt.Println(finishedTask)

		sendReportUrl := "http://" + Master + "/report"
		_, err := http.PostForm(sendReportUrl,
			url.Values{
				"fc": {string(cities)},
			})
		if err != nil {
			fmt.Println("Slave爬取情况发送异常" + err.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		finishedTask = finishedTask[:0]
		time.Sleep(time.Duration(5) * time.Second)
	}
	// defer resp.Body.Close()
}

func localIP() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)
	}
	return ips[1].String()
	// var localIP string
	// for _, ip := range ips {
	// 	localIP = ip.String()
	// 	if strings.Contains(localIP), ".") {
	// 		break
	// 	}
	// }

	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for _, address := range addrs {
	// 	// 检查IP地址，其他类型的地址(如link-local或者loopback)忽略
	// 	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
	// 		fmt.Println(ipnet.IP.String())
	// 	}
	// }

	//    import (
	//        "fmt"
	//        "github.com/tal-tech/go-zero/core/netx"
	//    )

	//    func main() {
	//        ip := netx.InternalIP()
	//        fmt.Println(ip)
	//    }
}
