package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type AuthUsecase interface {
	SSO(token string, authenticator SSOAuthenticator) (*entity.User, error)
	GenerateAuthToken(user entity.User, tokenizer Tokenizer) (string, error)
}

type DiscussionUsecase interface {
	GetAllUserDiscussions(userID string) ([]*entity.Discussion, error)
	GetDiscussionMessages(dicussionID string) ([]*entity.Message, error)
	Create(param entity.CreateDiscussionParam) (*entity.Discussion, error)
	SendMessage(dicussionID string, message entity.CreateMessage) (*entity.Message, error)
	Update(dicussionID string, param entity.UpdateDiscussionParam) (bool, error)
	Delete(dicussionID string) (bool, error)
}
