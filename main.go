package ratelimiter

import (
	"errors"
	"fmt"
	"time"

	deque "github.com/leexikang/generic-deque"
)

type RequestTimeStamps struct {
	timeStamps      deque.Deque[int64]
	Limit           int
	WindowTimeInSec int
}

func NewRequestTimeStamps(Limit, WindowTimeInSec int) *RequestTimeStamps {
	return &RequestTimeStamps{
    Limit: Limit,
    WindowTimeInSec: WindowTimeInSec,
		timeStamps: deque.Deque[int64]{},
	}
}

func (requestTimeStamps *RequestTimeStamps) Size() int {
	return requestTimeStamps.timeStamps.Len()
}

func (requestTimeStamps *RequestTimeStamps) isExceed() bool{
	return requestTimeStamps.timeStamps.Len()  >  requestTimeStamps.Limit
}

func (requestTimeStamps *RequestTimeStamps) Append(timeStamp int64) {
	requestTimeStamps.timeStamps.PushBack(timeStamp)
}

func (requestTimeStamps *RequestTimeStamps) EvictBefore(currentTime int64) {
	for requestTimeStamps.Size() != 0 &&
		requestTimeStamps.timeStamps.Front() < currentTime {
		requestTimeStamps.timeStamps.PopFront()
	}
}

type RateLimiter struct {
  requestTimeStamps map[string]*RequestTimeStamps 
}

func NewRateLimiter() *RateLimiter{
  return &RateLimiter{
    requestTimeStamps: make(map[string]*RequestTimeStamps),
  }
}

func (ratelimiter *RateLimiter) create(id string, timestamps RequestTimeStamps) {
  ratelimiter.requestTimeStamps[id] = &timestamps
}


func (ratelimiter *RateLimiter) insert(id string) error{
  requestTimeStamps, ok := ratelimiter.requestTimeStamps[id] 
  if !ok {
    panic("user has not been initialized yet")
  }

  currentTime := time.Now().UnixMilli()
  requestTimeStamps.Append(currentTime)
  requestTimeStamps.EvictBefore(currentTime - int64(requestTimeStamps.WindowTimeInSec * 1000))

  if requestTimeStamps.isExceed() {
  errMessage := fmt.Sprintf("Your are exceed than the limit of %d in %d seconds", requestTimeStamps.Limit, requestTimeStamps.WindowTimeInSec)
    return errors.New(errMessage)
  }

  return nil
}
