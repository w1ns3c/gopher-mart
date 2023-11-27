package middlewares

import (
	"net/http"
	"sync"
	"time"
)

type ddosMiddleware struct {
	maxRequestNum uint64
	currentNum    uint64
	startTime     time.Time
	mu            *sync.Mutex
}

var (
	defaultTime = time.Date(2000, 01, 01, 01, 01, 01, 01, time.Local)
)

func NewDDOSMiddleware(maxRequestNum uint64) *ddosMiddleware {
	return &ddosMiddleware{
		maxRequestNum: maxRequestNum,
		currentNum:    0,
		startTime:     defaultTime,
		mu:            &sync.Mutex{},
	}
}

func (m *ddosMiddleware) DDOSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m.mu.Lock()
		defer m.mu.Unlock()
		// if first request
		if m.startTime == defaultTime {
			m.startTime = time.Now()
		}

		// if 1 minute has passed
		if time.Since(m.startTime) > time.Minute {
			m.currentNum = 0
		}

		m.currentNum += 1

		// if too many requests
		if m.currentNum > m.maxRequestNum {
			writer.Header().Add("retry-after", "60")
			writer.WriteHeader(http.StatusTooManyRequests)
			return
		}
	})
}
