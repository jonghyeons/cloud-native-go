package ch04

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	var count int
	var EmulateTransientError = func(ctx context.Context) (string, error) {
		count++
		if count <= 3 {
			return "intentional fail", errors.New("error")
		} else {
			return "success", nil
		}
	}

	r := Retry(EmulateTransientError, 5, 2*time.Second)
	res, err := r(context.Background())

	fmt.Println(res, err)
}
