package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	userMetaPrefix   = "user_limit"
	timestampsPrefix = "user_timestamps_"
	windowInSecKey   = "windowInSec"
	limitKey         = "limit"
)

func main() {
	redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type RateLimiter struct {
	client redis.Client
}

func (ratelimiter *RateLimiter) create(ctx context.Context, id string, limit, windowInSec int) {
	ratelimiter.client.HSet(ctx, userMetaPrefix+id, "limit", limit, "windowInSec", windowInSec)
}

func (ratelimiter *RateLimiter) delte(ctx context.Context, id string, limit, windowInSec int) {
	ratelimiter.client.HDel(ctx, userMetaPrefix+id)
	ratelimiter.client.ZRem(ctx, timestampsPrefix+id)
}

func (ratelimiter *RateLimiter) isAllowed(ctx context.Context, id string) error {
	limitInStr, err := ratelimiter.client.HGet(ctx, userMetaPrefix+id, "limit").Result()
	if err != nil {
		return errors.New("User has not been initialized")
	}

	var limit int64
	if limit, err = strconv.ParseInt(limitInStr, 10, 64); err != nil {
		return errors.New("Error to parse windowInSec")
	}

	windowInSecStr, err := ratelimiter.client.HGet(ctx, userMetaPrefix+id, windowInSecKey).Result()
	if err != nil {
		return errors.New("User has not been initialized")
	}

	var windowInSec int64
	if windowInSec, err = strconv.ParseInt(windowInSecStr, 10, 64); err != nil {
		return errors.New("Error to parse windowInSec")
	}

	currentTime := time.Now().UnixMilli()
	startOfWindow := strconv.FormatInt((currentTime - (windowInSec * 1000)), 10)
	_, err = ratelimiter.client.ZRemRangeByScore(ctx, timestampsPrefix+id, "0", startOfWindow).
		Result()
	if err != nil {
		return err
	}

	member := redis.Z{
		Score:  float64(currentTime),
		Member: currentTime,
	}
	_, err = ratelimiter.client.ZAdd(ctx, timestampsPrefix+id, &member).Result()
	if err != nil {
		return err
	}

	count, err := ratelimiter.client.ZCount(ctx, timestampsPrefix+id, "0", "9999999999999").Result()
	if err != nil {
		return err
	}

	if count > limit {
		return errors.New("Limited exceed")
	}

	return nil
}
