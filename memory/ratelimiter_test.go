package memory

import (
	"strconv"
	"testing"
	"time"
)

func TestAppendTimeStamps(t *testing.T) {
	request := newRequestTimeStamps(10, 10)
	currentTime := time.Now().Unix()
	times := []int64{
		currentTime,
		currentTime,
		currentTime,
	}

	for i := 0; i < len(times); i++ {
		request.Append(times[i])
	}

	if request.Size() != len(times) {
		t.Errorf("Expected size to be %d but found %d", len(times), request.Size())
	}
}

func TestEvictBefore(t *testing.T) {
	request := newRequestTimeStamps(10, 10)
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

type Request struct {
	ID int64
}

func (request Request) Id() string {
	return strconv.FormatInt(request.ID, 10)
}

func TestRateLimiter(t *testing.T) {
	limit := 2
	windowInSec := 1
	limiter := New()
	id := "1"
	limiter.create(id, limit, windowInSec)
	limiter.insert(id)
	limiter.insert(id)
	err := limiter.insert(id)

	if err == nil {
		t.Errorf("Limit executed %d but don't throw error", limit)
	}

	limiter.delete(id)

	err = limiter.insert(id)

	if err == nil {
		t.Errorf("Error not thrown, when inserting with delted id")
	}
}
