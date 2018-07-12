package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer res.Body.Close()
	fmt.Println("Response received, status code:", res.StatusCode)
}
