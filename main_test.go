package ratelimiter

import (
	"testing"
	"time"
)

func TestAppendTimeStamps(t *testing.T) {
	request := NewRequestTimeStamps(10, 10)
	currentTime := time.Now().Unix()
	times := []int64{
		currentTime,
		currentTime + 1,
		currentTime + 2,
	}

	for i := 0; i < len(times); i++ {
		request.Append(times[i])
	}

	if request.Size() != len(times) {
		t.Errorf("Expected size to be %d but found %d", len(times), request.Size())
	}
}

func TestEvictBefore(t *testing.T) {
	request := NewRequestTimeStamps(10, 10)
	currentTime := int64(0)
	times := []int64{
		currentTime,
		currentTime + 1,
		currentTime + 2,
	}

	for i := 0; i < len(times); i++ {
		request.Append(times[i])
	}

	request.EvictBefore(currentTime + 1)

	if request.Size() != 2 {
		t.Errorf("Expected size to be %d but found %d", 2, request.Size())
	}

	if request.timeStamps.Front() != times[1] {
		t.Errorf("Expected timeStamp to be %d but found %d", times[1], request.timeStamps.Front())
	}

	request.EvictBefore(currentTime + 2)

	if request.Size() != 1 {
		t.Errorf("Expected size to be %d but found %d", 1, request.Size())
	}

	if request.timeStamps.Front() != times[2] {
		t.Errorf("Expected timeStamp to be %d but found %d", times[2], request.timeStamps.Front())
	}

	request.EvictBefore(currentTime + 3)

	if request.Size() != 0 {
		t.Errorf("Expected size to be %d but found %d", 0, request.Size())
	}
}
