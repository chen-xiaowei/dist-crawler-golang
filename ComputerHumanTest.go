package main

import (
	"fmt"
	// "fmt"

	"github.com/tebeka/selenium/chrome"

	// "os"
	// "strings"
	"time"

	"github.com/tebeka/selenium"
)

func main1() {
	ops := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService("D:\\workspace\\enterprise-searcher\\driver\\chromedriver.exe", 8088, ops...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}

	chromeCaps := chrome.Capabilities{
		ExcludeSwitches: []string{"enable-automation"},
		Args:            []string{"no-sandbox"},
	}
	caps.AddChrome(chromeCaps)
	// wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 8088))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	url := "http://search.fang.com/captcha-403881c98f414fc78f/?t=1689819078.893&h=aHR0cHM6Ly95ei5lc2YuZmFuZy5jb20vaG91c2UtYTAxNjgxMT9yZnNzPTEtNDAzODgxYzk4ZjQxNGZjNzhmLTAx&c=cmE6MTE4LjI0Ny4xMTcuNDE7eHJpOjt4ZmY6"
	if err := wd.Get(url); err != nil {
		panic(err)
	}
	//drag-handler verifyicon center-icon
	ele, _ := wd.FindElement(selenium.ByClassName, "drag-handler")
	fmt.Println(ele.GetAttribute("class"))
	// wd.KeyDown()
	p, _ := ele.LocationInView()
	fmt.Printf("%d %d \n", p.X, p.Y)

	// args := make([]interface{}, 0)
	// args = append(args, ele)
	// js := "var dragBtn = document.getElementsByClassName('drag-handler')[0];"
	// js += "var dragImg = document.getElementsByClassName('img-block')[0];"
	// js += "dragBtn.style.left = dragBtn.offsetLeft + 50 + 'px';"
	// js += "dragImg.style.left = dragImg.offsetLeft + 50 + 'px';"
	// js += "dragBtn.click();"
	// js += "dragBtn.onmouseup();"
	// .dispatchEvent(new MouseEvent('click', {shiftKey: true}))
	// res, err := wd.ExecuteScript(js, args)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)

	time.Sleep(2 * time.Second)
	if err := ele.MoveTo(p.X, p.Y); err != nil {
		panic(err)
	}

	if err = wd.KeyDown("mouse left"); err != nil {
		panic(err)
	}
	var selected bool
	if selected, err = ele.IsSelected(); err != nil {
		panic(err)
	}
	fmt.Println(selected)

	time.Sleep(2 * time.Second)
	if err = wd.ButtonUp(); err != nil {
		panic(err)
	}

	// if err = ele.Click(); err != nil {
	// 	panic(err)
	// }
	time.Sleep(50 * time.Second)
}
