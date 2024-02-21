package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Define the struct wrapper around raw redis client
type StoreService struct {
	redisClient *redis.Client
}

// Define the global variables for the store service and context
var (
	storeService = &StoreService{}
	ctx          = context.Background()
)

// Initialize the store service with the redis client
func InitializeStore() *StoreService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Printf("Connected to Redis: %v\n", pong)
	storeService.redisClient = redisClient
	return storeService
}

// Map the short URL to the original URL
func SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed to save URL mapping: %v", err))
	}
}

// Retrieve the original URL from the short URL
func RetrieveOriginalUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to retrieve original URL: %v", err))
	}
	return result
}

// Note that in a real world usage, the cache duration shouldn't have
// an expiration time, an LRU policy config should be set where the
// values that are retrieved less often are purged automatically from
// the cache and stored back in RDBMS whenever the cache is full
const CacheDuration = 6 * time.Hour
