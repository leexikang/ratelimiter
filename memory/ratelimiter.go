package memory

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	deque "github.com/leexikang/generic-deque"
)

type requestTimeStamps struct {
	mu              *sync.Mutex
	timeStamps      deque.Deque[int64]
	limit           int
	windowTimeInSec int
}

func newRequestTimeStamps(Limit, WindowTimeInSec int) *requestTimeStamps {
	return &requestTimeStamps{
		limit:           Limit,
		mu:              &sync.Mutex{},
		windowTimeInSec: WindowTimeInSec,
		timeStamps:      deque.Deque[int64]{},
	}
}

func (requestTimeStamps *requestTimeStamps) Size() int {
	return requestTimeStamps.timeStamps.Len()
}

func (requestTimeStamps *requestTimeStamps) isExceed() bool {
	return requestTimeStamps.timeStamps.Len() > requestTimeStamps.limit
}

func (requestTimeStamps *requestTimeStamps) Append(timeStamp int64) {
	requestTimeStamps.mu.Lock()
	requestTimeStamps.timeStamps.PushBack(timeStamp)
	requestTimeStamps.mu.Unlock()
}

func (requestTimeStamps *requestTimeStamps) EvictBefore(currentTime int64) {
	requestTimeStamps.mu.Lock()
	for requestTimeStamps.Size() != 0 &&
		requestTimeStamps.timeStamps.Front() < currentTime {
		requestTimeStamps.timeStamps.PopFront()
	}
	requestTimeStamps.mu.Unlock()
}

type RateLimiter struct {
	requestTimeStamps map[string]*requestTimeStamps
}

func New() *RateLimiter {
	return &RateLimiter{
		requestTimeStamps: make(map[string]*requestTimeStamps),
	}
}

func (ratelimiter *RateLimiter) Create(id string, limit, windowTimeInSec int) {
	ratelimiter.requestTimeStamps[id] = newRequestTimeStamps(limit, windowTimeInSec)
}

func (ratelimiter *RateLimiter) Delete(id string) {
	delete(ratelimiter.requestTimeStamps, id)
}

func (ratelimiter *RateLimiter) IsAllowed(id string) error {
	requestTimeStamps, ok := ratelimiter.requestTimeStamps[id]

	if !ok {
		return errors.New("user has not been initialized yet")
	}

	currentTime := time.Now().UnixMilli()
	log.Print(requestTimeStamps.windowTimeInSec)
	requestTimeStamps.EvictBefore(currentTime - int64(requestTimeStamps.windowTimeInSec*1000))
	requestTimeStamps.Append(currentTime)
	if requestTimeStamps.isExceed() {
		errMessage := fmt.Sprintf(
			"Your are exceed than the limit of %d in %d seconds",
			requestTimeStamps.limit,
			requestTimeStamps.windowTimeInSec,
		)
		return errors.New(errMessage)
	}

	return nil
}
