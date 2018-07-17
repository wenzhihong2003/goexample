package main

import (
	"reflect"
)

type inf interface {
	Method1()
	Method2()
}

type ss struct {
	a func()
}

func (i ss) Method1() {

}

func (i ss) Method2() {

}

func main() {
	s := reflect.TypeOf(ss{})
	i := reflect.TypeOf(new(inf)).Elem()

}
