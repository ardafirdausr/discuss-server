package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	Create(param entity.CreateUserParam) (*entity.User, error)
}

type MessageRepository interface {
	Create(param entity.CreateMessage) (*entity.Message, error)
	GetMessagesByDiscussionID(discussionID string) ([]*entity.Message, error)
}

type DiscussionRepository interface {
	GetDissionsByUserID(userID string) ([]*entity.Discussion, error)
	Create(param entity.CreateDiscussionParam) (*entity.Discussion, error)
	UpdateByID(ID string, param entity.UpdateDiscussionParam) (bool, error)
	DeleteByID(discussionID string) (bool, error)
}

type PubSubRepository interface {
	Publish(channel string, message string) error
	Subscribe(channels ...string) (<-chan interface{}, error)
	Unsubscribe(channels ...string) error
}
