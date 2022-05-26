package ratelimiter

type RateLimiter interface {
	create(string)
	delet()
	insert(string) error
}
