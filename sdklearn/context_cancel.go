package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// 在 Go 中用 Context 取消操作
// https://github.com/studygolang/GCTT/blob/master/translated/tech/20180617-Using-context-cancellation-in-Go.md

func main() {
	c1()
}

// 监听取消事件
// context 包提供了 Done() 方法, 它返回一个当 Context 收取到 取消 事件时会接收到一个 struct{} 类型的 channel。 监听取消事件只需要简单的等待 <- ctx.Done() 就好了例如： 一个 Http Server 会花2秒去处理事务，如果请求提前取消，我们想立马返回结果：
// 你可以通过执行这段代码, 用浏览器打开 localhost:8000。如果你在2秒内关闭浏览器，你会看到在控制台打印了 "request canceled"。
func c1() {
	// create an http server that listens on port 8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "processing request\n")
		// We use `select` to execute a peice of code depending on which
		// channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// If we receive a message after 2 seconds
			// that means the request has been processed
			// We then write this as the response
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	}))
}

// 触发取消事件
// 如果你有一个可以取消的操作，你可以通过context触发一个 取消事件 。 这个你可以用 context 包 提供的 WithCancel 方法， 它返回一个 context 对象，和一个没有参数的方法。这个方法不会返回任何东西，仅在你想取消这个context的时候去调用。
//
// 第二种情况是依赖。 依赖的意思是，当一个操作失败，会导致其他操作失败。 例如：我们提前知道了一个操作失败，我们会取消所有依赖操作。

func operation1(ctx context.Context) error {
	// Let's assume that this operation failed for some reason
	// We use time.Sleep to simulate a resource intensive operation
	time.Sleep(100 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	// We use a similar pattern to the HTTP server
	// that we saw in the earlier example
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halted operation2")
	}
}

func c2() {
	// Create a new context
	ctx := context.Background()
	// Create a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(ctx)

	// Run two operations: one in a different go routine
	go func() {
		err := operation1(ctx)
		// If this operation returns an error
		// cancel all operations using this context
		if err != nil {
			cancel()
		}
	}()

	// Run operation2 with the same context we use for operation1
	operation2(ctx)
}

// 基于时间的取消操作
// 任何程序对一个请求的最大处理时间都需要维护一个 SLA (服务级别协议)，这可以使用基于时间的取消。这个 API 基本和上一个例子相同，只是多了一点点：
// 例如:HTTP API 调用内部的服务。 如果服务长时间没有响应，最好提前返回失败并取消请求。
func c3() {
	// Create a new context
	// With a deadline of 100 milliseconds
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	// Make a request, that will call the google homepage
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	// Associate the cancellable context we just created to the request
	req = req.WithContext(ctx)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	// If the request failed, log to STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	// Print the statuscode if the request succeeds
	fmt.Println("Response received, status code:", res.StatusCode)
}

// 陷阱和注意事项
// 尽管 Go 的 context 很好用，但是在使用之前最好记住几点。最重要的就是：context 仅可以取消一次。如果想传播多个错误的话，context 取消并不是最好的选择，最惯用的场景是你真的想取消一个操作，并通知下游操作发生了一个错误。
//
// 另一个要记住的是，一个 context 实例会贯穿所有你想使用取消操作的方法和 go-routines 。要避免使用一个已取消的 context 作为 WithTimeout 或者 WithCancel 的参数，这可能导致不确定的事情发生。
