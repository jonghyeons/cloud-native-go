package ch04

import (
	"context"
	"errors"
	"sync"
	"time"
)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailure int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := consecutiveFailure - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}
		m.RUnlock()
		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()
		if err != nil {
			consecutiveFailure++
			return response, err
		}

		consecutiveFailure = 0
		return response, nil
	}
}
