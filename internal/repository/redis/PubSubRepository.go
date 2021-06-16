package redis

import "github.com/go-redis/redis"

type PubSubRepository struct {
	client *redis.Client
}

func NewPubSubRepository(client *redis.Client) *PubSubRepository {
	return &PubSubRepository{client: client}
}

func (rr PubSubRepository) Publish(channel string, message string) error {
	return nil
}

func (rr PubSubRepository) Subscribe(channels ...string) (<-chan interface{}, error) {
	return nil, nil
}

func (rr PubSubRepository) Unsubscribe(channels ...string) error {
	return nil
}
