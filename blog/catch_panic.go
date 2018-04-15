package main

import (
	"runtime"
	"bytes"
)

// Golang捕获panic堆栈信息的优雅姿势
// https://my.oschina.net/henrylee2cn/blog/873885

// 该函数的优点：
//
// 1. 比直接recover()捕获的panic信息更加详尽
// 2. 比直接放任其panic打印的堆栈信息更精准，第一行就是发生panic的代码行
// 3. 比直接放任其panic打印的堆栈信息更简洁，可以指定信息量（kb）

// PanicTrace trace panic stack info
func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) // 4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}

func main() {

}
