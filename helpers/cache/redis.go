package cache

import "github.com/go-redis/redis"

func New() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return
}

func DefaultClient() (client *redis.Client) {
	return New()
}
