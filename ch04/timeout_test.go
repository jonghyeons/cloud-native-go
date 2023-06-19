package ch04

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var Slow = func(s string) (string, error) {
		if s == "" {
			return "", errors.New("string parameter is empty")
		}
		return s, nil
	}
	timeout := Timeout(Slow)
	res, err := timeout(ctxt, "some input")

	fmt.Println(res, err)
}
