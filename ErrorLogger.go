package main

import (
	"bufio"
	"fmt"
	"os"
)

var writer *bufio.Writer

func init() {
	file, err := os.OpenFile("./errlog.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("日志文件打开异常= %v \n", err)
	}
	writer = bufio.NewWriter(file)
}

func Err(err string) {
	fmt.Println(err)
	if !recordLog(err) {
		writer.WriteString(err + "\r\n")
		writer.Flush()
	}
}
