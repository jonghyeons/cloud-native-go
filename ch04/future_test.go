package ch04

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	var slot = func(ctx context.Context) Future {
		resCh := make(chan string)
		errCh := make(chan error)

		go func() {
			select {
			case <-time.After(time.Second * 2):
				resCh <- "I slept for 2 seconds"
				errCh <- nil
			case <-ctx.Done():
				resCh <- ""
				errCh <- ctx.Err()
			}
		}()

		return &InnerFuture{resCh: resCh, errCh: errCh}
	}

	ctx := context.Background()
	future := slot(ctx)

	res, err := future.Result()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(res)
}
