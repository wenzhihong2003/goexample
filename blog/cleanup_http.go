package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 参考: https://gianarb.it/blog/go-http-cleanup-http-connection-terminated
// 清除在http长时间运行时, 用户断开连接后, 进行资源清除

func main() {
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		err := ioutil.WriteFile(os.TempDir()+"/txt", []byte("hello"), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("new file " + os.TempDir() + "/txt")
		notify := w.(http.CloseNotifier).CloseNotify()  // 重点在这里
		go func() {
			<-notify  // 用户过早中断连接时, 会得到一个消息
			fmt.Println("the client closed the connection prematurely. Cleaning up.")
			os.Remove(os.TempDir() + "/txt")
		}()
		time.Sleep(4 * time.Second)
		fmt.Fprintln(w, "File persisted.")
	})
	http.ListenAndServe(":8080", nil)
}
