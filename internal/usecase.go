package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type AuthUsecase interface {
	SSO(token string, authenticator SSOAuthenticator) (*entity.User, error)
	GenerateAuthToken(user entity.User, tokenizer Tokenizer) (string, error)
}

type DiscussionUsecase interface {
	GetAllUserDiscussions(userID interface{}) ([]*entity.Discussion, error)
	GetDiscussionByID(ID interface{}) (*entity.Discussion, error)
	GetDiscussionByCode(code string) (*entity.Discussion, error)
	Create(param entity.CreateDiscussionParam) (*entity.Discussion, error)
	Update(discussionID interface{}, param entity.UpdateDiscussionParam) error
	UpdatePassword(discussionID interface{}, param entity.UpdateDiscussionPassword) error
	UpdatePhoto(discussionID interface{}, url string) error
	Delete(discussionID interface{}) error
}

type MessageUsecase interface {
	SendMessage(pubsub PubSub, message entity.CreateMessage) (*entity.Message, error)
}
