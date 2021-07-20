package pubsub

import (
	"log"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/go-redis/redis"
)

type RedisPubSub struct {
	client *redis.Client
	pubsub *redis.PubSub
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	pubsub := client.Subscribe()
	return &RedisPubSub{
		client: client,
		pubsub: pubsub,
	}
}

func (rps RedisPubSub) Publish(channel string, message interface{}) error {
	icmd := rps.client.Publish(channel, message)
	return icmd.Err()
}

func (rps RedisPubSub) Subscribe(channels ...string) error {
	if err := rps.pubsub.Subscribe(channels...); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (rps RedisPubSub) Unsubscribe(channels ...string) error {
	if err := rps.pubsub.Unsubscribe(channels...); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (rps RedisPubSub) Listen(listener internal.SubscribeListener) {
	rdsMessages := rps.pubsub.Channel()
	for rdsMessage := range rdsMessages {
		listener(rdsMessage.Channel, rdsMessage.Payload)
	}
}

func (rps RedisPubSub) Close() error {
	return rps.pubsub.Close()
}
