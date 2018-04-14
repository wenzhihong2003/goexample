package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"fmt"
)

var engine *xorm.Engine

func main() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@127.0.0.1:3307/apiserv?charset=utf8")
	if err != nil {
		panic(err)
	}else {
		fmt.Println("ok...........")
	}
}
