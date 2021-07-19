package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type UserRepository interface {
	GetUserByID(userID interface{}) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	Create(param entity.CreateUserParam) (*entity.User, error)
}

type MessageRepository interface {
	Create(param entity.CreateMessage) (*entity.Message, error)
	GetMessagesByDiscussionID(discussionID interface{}) ([]*entity.Message, error)
}

type DiscussionRepository interface {
	GetDiscussionsByID(discussionID interface{}) (*entity.Discussion, error)
	GetDiscussionByCode(code string) (*entity.Discussion, error)
	GetDiscussionsByUserID(userID interface{}) ([]*entity.Discussion, error)
	Create(param entity.CreateDiscussionParam) (*entity.Discussion, error)
	AddMember(code string, userID interface{}) error
	RemoveMember(discussionID interface{}, userID interface{}) error
	UpdateByID(discussionID interface{}, param entity.UpdateDiscussionParam) error
	UpdatePasswordByID(discussionID interface{}, password string) error
	UpdatePhotoByID(discussionID interface{}, url string) error
	DeleteByID(discussionID interface{}) error
}
