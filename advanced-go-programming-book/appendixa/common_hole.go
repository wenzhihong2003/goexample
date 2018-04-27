package main

import (
	"fmt"
)

// 附录A：Go语言常见坑
// https://chai2010.gitbooks.io/advanced-go-programming-book/content/appendix/appendix-a-trap.html

var msg string
var done bool = false

func main() {

}

// 数组是值传递
// 在函数调用参数中, 数组是值传递, 无法通过修改数组类型的参数返回结果.
// 必要时需要使用切片.

func arrUseParamVal()  {
	x := [3]int{1, 2, 3}
	func(arr [3]int){
		arr[0] = 7
		fmt.Println(arr)
	}(x)

	fmt.Println(x)
}
