package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Define the struct wrapper around raw Redis client
type StoreService struct  {
	redisClient *redis.Client
}

// Top level declarations for the storeService and Redis context
var (
	storageService = &StoreService{}
	ctx = context.Background()
)

const CasheDuration = 6 * time.Hour



/* initliazation connection for redis
*/
func InitializationRedisClient() *StoreService {
	port := 6379
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%v", port),
		Password: "",
		DB: 0,
	})

	// in older verision of redis ping don't take a ctx in ping 
	// so if you want to adding this you must update the redis to v8 or v9
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error happend in redis connection: %v", err)
	}

	fmt.Printf("redis Client running in port: %v, pong message: %v", port, pong)
	storageService.redisClient = client
	return storageService
}


/* We want to be able to save the mapping between the originalUrl 
and the generated shortUrl url
*/
func SaveUrlMapping(shortUrl string, originalUrl string, userId string){ 
	err := storageService.redisClient.Set(ctx, shortUrl, originalUrl, 
		CasheDuration).Err()
	if err != nil {
		log.Fatalf("error happend in saving data redis: %v", err)
	}
}

/* We should be able to retrieve the initial long URL once the short 
is provided. This is when users will be calling the shortlink in the 
url, so what we need to do here is to retrieve the long url and
think about redirect.
*/
func RetrieveInitialUrl(shortUrl string) string {
	result, err := storageService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		log.Fatalf("error happend in retreiving data redis: %v", err)
	}
	return result
}
