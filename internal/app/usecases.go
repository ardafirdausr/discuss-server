package app

import (
	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/usecase"
)

type UseCases struct {
	AuthUsecase internal.AuthUsecase
}

func newUseCases(repos *Repositories) *UseCases {
	authUsecase := usecase.NewAuthUsecase(repos.userRepo)
	return &UseCases{
		AuthUsecase: authUsecase,
	}
}
