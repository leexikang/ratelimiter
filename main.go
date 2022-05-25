package ratelimiter

import (

	deque "github.com/leexikang/generic-deque"
)

type RequestTimeStamps struct {
	timeStamps      deque.Deque[int64]
	Limit           int
	WindowTimeInSec int
}

func NewRequestTimeStamps(Limit, WindowTimeInSec int) *RequestTimeStamps {
	return &RequestTimeStamps{
		timeStamps: deque.Deque[int64]{},
	}
}

func (requestTimeStamps *RequestTimeStamps) Size() int {
	return requestTimeStamps.timeStamps.Len()
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


// Create a map with UserId key and the value is requestTimeStamps
//  If user is not created create a new one for it
// Provide append method
//  Append the timeStamp to the requestTimeStamps 
//  Evict the logs with timestamp less than now - WindowTimeInSec 
//  compare the logs size with the RequestTimeStamps' Limit, 
//  if greateer allow, if less than denied

type RateLimiter struct {
  requestTimeStamps map[string]RequestTimeStamps 
}

type request interface {
  Id() string 
}

func (ratelimiter *RateLimiter) insert(request request) {
  requestTimeStamps, ok := ratelimiter.requestTimeStamps[request.Id()] 
  if !ok {
    panic("user has not been initialized yet")
  }
}
