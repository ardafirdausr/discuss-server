package internal

import (
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type SSOAuthenticator interface {
	Authenticate(token string) (*entity.User, error)
}

type Tokenizer interface {
	Generate(entity.TokenPayload) (string, error)
}

type PubSub interface {
	Publish(channel string, message interface{}) error
	Subscribe(channels ...string) error
	Unsubscribe(channels ...string) error
	Listen(listener SubscribeListener)
	Close() error
}

type SubscribeListener func(channel string, message interface{}) error
