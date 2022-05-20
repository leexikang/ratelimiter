package ratelimiter

import deque "github.com/leexikang/generic-deque"

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
