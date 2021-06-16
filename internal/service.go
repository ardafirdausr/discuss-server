package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type SSOAuthenticator interface {
	Authenticate(token string) (*entity.User, error)
}

type Tokenizer interface {
	Generate(entity.TokenPayload) (string, error)
}

type Mailer interface {
	SendMail(entity.Mail) error
}

type PubSub interface {
}
