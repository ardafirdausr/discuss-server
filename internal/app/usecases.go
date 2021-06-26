package app

import (
	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/usecase"
)

type Usecases struct {
	AuthUsecase       internal.AuthUsecase
	DiscussionUsecase internal.DiscussionUsecase
}

func newUsecases(repos *Repositories) *Usecases {
	authUsecase := usecase.NewAuthUsecase(repos.userRepo)
	discussionUsecase := usecase.NewDiscussionUsecase(repos.discussionRepo, repos.messageRepo)
	return &Usecases{
		AuthUsecase:       authUsecase,
		DiscussionUsecase: discussionUsecase,
	}
}
