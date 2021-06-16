package app

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/ardafirdausr/discuss-server/internal"
	mongoRepo "github.com/ardafirdausr/discuss-server/internal/repository/mongo"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	userRepo internal.UserRepository
}

func newRepositories() (*Repositories, error) {
	mongoDB, err := connectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database\n%v", err)
		return nil, err
	}

	// redis, err := connnectToRedis()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to the database\n%v", err)
	// 	return nil, err
	// }

	userRepo := mongoRepo.NewUserRepository(mongoDB)
	// pubsubRepo := redisRepo.NewPubSubRepository(redis)

	return &Repositories{
		userRepo: userRepo,
	}, nil
}

func connectToMongoDB() (*mongo.Database, error) {
	mongoDBURI := os.Getenv("MONGO_DB_URI")
	DBName := os.Getenv("MONGO_DB_NAME")

	clientOptions := options.Client().ApplyURI(mongoDBURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database(DBName), nil
}

func connnectToRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDRESS")
	pwd := os.Getenv("REDIS_PASSWORD")

	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
	rdsOpt := &redis.Options{
		Addr:      addr,
		Password:  pwd,
		TLSConfig: tlsCfg,
	}

	rdsClient := redis.NewClient(rdsOpt)
	status := rdsClient.Ping()
	if err := status.Err(); err != nil {
		return nil, err
	}

	return rdsClient, nil
}
