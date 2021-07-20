package app

import (
	"github.com/ardafirdausr/discuss-server/internal/driver"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type drivers struct {
	Mongo *mongo.Database
	Redis *redis.Client
}

func newDrivers() (*drivers, error) {
	drivers := &drivers{}

	mongo, err := driver.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	drivers.Mongo = mongo

	redis, err := driver.ConnnectToRedis()
	if err != nil {
		return nil, err
	}
	drivers.Redis = redis

	return drivers, nil
}
