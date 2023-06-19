package ch04

import (
	"fmt"
	"testing"
	"time"
)

func TestFanIn(t *testing.T) {
	sources := make([]<-chan int, 0) // 빈 채널 슬라이스를 생성합니다

	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)

		go func() {
			defer close(ch)

			for i := 1; i <= 5; i++ {
				ch <- i
				time.Sleep(time.Second)
			}
		}()

		dest := Funnel(sources...)
		for d := range dest {
			fmt.Println(d)
		}
	}
}
