package redis

import (
	"context"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	conf := RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	rateLimiter := New(conf)
	ctx := context.Background()
	id := "1"
	limit := 2
	windowInSec := 1
	rateLimiter.create(ctx, id, limit, windowInSec)
	err := rateLimiter.isAllowed(ctx, id)
	if err != nil {
		t.Error(err)
	}

	err = rateLimiter.isAllowed(ctx, id)
	if err != nil {
		t.Error(err)
	}

	err = rateLimiter.isAllowed(ctx, id)
	if err == nil {
		t.Error("should return limit exceeded error")
	}
}
