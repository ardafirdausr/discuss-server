package driver

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Database, error) {
	mongoDBURI := os.Getenv("MONGO_DB_URI")
	DBName := os.Getenv("MONGO_DB_NAME")

	clientOptions := options.Client().ApplyURI(mongoDBURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return client.Database(DBName), nil
}
