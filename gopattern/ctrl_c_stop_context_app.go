package main

import (
	"context"
	"os"
	"os/signal"
)

// Make Ctrl+C cancel the context.Context
// 参考: https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff

func main() {
	ctx := context.Background()

	// trap ctrl+c and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	// doSomethinAwesom(ctx)
}
