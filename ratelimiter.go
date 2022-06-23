package ratelimiter

type RateLimiter interface {
	Create(id string, limit, windowTimeInSec int)
	Delete(string)
	IsAllowed(string) error
}
