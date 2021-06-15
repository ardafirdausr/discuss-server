package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type SSOAuthenticator interface {
	Authenticate()
}

type Tokenizer interface {
	Generate(entity.TokenPayload) (string, error)
}
