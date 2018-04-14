package main

import (
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func main() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@127.0.0.1:3307/apiserv?charset=utf8")
	if err != nil {
		panic(err)
	}
}
