package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	Create(param entity.CreateUserParam) (*entity.User, error)
}

type MessageRepository interface {
	Create(param entity.CreateMessage) (*entity.Message, error)
	GetMessagesByDiscussionID(discussionID interface{}) ([]*entity.Message, error)
}

type DiscussionRepository interface {
	GetDiscussionsByID(ID interface{}) (*entity.Discussion, error)
	GetDiscussionsByCode(code string) (*entity.Discussion, error)
	GetDiscussionsByUserID(userID interface{}) ([]*entity.Discussion, error)
	Create(param entity.CreateDiscussionParam) (*entity.Discussion, error)
	UpdateByID(ID interface{}, param entity.UpdateDiscussionParam) error
	UpdatePasswordByID(ID interface{}, password string) error
	UpdatePhotoByID(ID interface{}, url string) error
	DeleteByID(discussionID interface{}) error
}

type PubSubRepository interface {
	Publish(channel string, message string) error
	Subscribe(channels ...string) (<-chan interface{}, error)
	Unsubscribe(channels ...string) error
}
