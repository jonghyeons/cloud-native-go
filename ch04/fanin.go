package ch04

import "sync"

func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int) // 공유 출력 채널 선언

	// 모든 source의 채널이 닫혔을 때 출력 채널을 자동으로 닫기 위해 사용됩니다.
	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, ch := range sources { // 각 입력 채널에 대해 고루틴을 시작합니다
		go func(c <-chan int) {
			defer wg.Done() // 채널이 닫히면 WaitGroup으로 알려줍니다.

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() { // 모든 입력 채널이 닫힌 후
		wg.Wait() // 출력 채널을 닫기 위한 고루틴을 시작합니다
		close(dest)
	}()

	return dest
}
