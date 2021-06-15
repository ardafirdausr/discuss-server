package internal

import "github.com/ardafirdausr/discuss-server/internal/entity"

type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	Create(param entity.CreateUserParam) (*entity.User, error)
}
