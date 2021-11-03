package main

import (
	"course/week_02_error"
	"course/week_03_errgroup"
)

// 第二周作业调用
func week_02() {
	week_02_error.GetUserNameById(1)
}

// 第三周作业调用
func week_03() {
	week_03_errgroup.HttpServer()
}

func main() {
	week_03()
}
