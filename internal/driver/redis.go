package driver

import (
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis"
)

func ConnnectToRedis() (*redis.Client, error) {
	redisURI := os.Getenv("REDIS_URI")
	redisUrl, err := url.Parse(redisURI)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	redisPassword, _ := redisUrl.User.Password()
	redisDB := 0
	redisOptions := &redis.Options{
		Addr:     redisUrl.Host,
		Password: redisPassword,
		DB:       redisDB,
	}

	rdsClient := redis.NewClient(redisOptions)
	status := rdsClient.Ping()
	if err := status.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return rdsClient, nil
}
