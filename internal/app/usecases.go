package app

import (
	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/usecase"
)

type Usecases struct {
	AuthUsecase internal.AuthUsecase
}

func newUsecases(repos *Repositories) *Usecases {
	authUsecase := usecase.NewAuthUsecase(repos.userRepo)
	return &Usecases{
		AuthUsecase: authUsecase,
	}
}
