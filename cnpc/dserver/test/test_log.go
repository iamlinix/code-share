package main

import (
	"cnpc.com.cn/cnpc/dserver/zaps"
)


func main() {

	zaps.Info("this is a info log message")
	zaps.Warn("this is a warn log message")
	zaps.Error("<<< this is an error with code 20005 ")
}
