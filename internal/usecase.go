package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type AuthUsecase interface {
	SSO(token string, authenticator SSOAuthenticator) (*entity.User, error)
	GenerateAuthToken(user entity.User, tokenizer Tokenizer) (string, error)
}
