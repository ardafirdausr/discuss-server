package app

import (
	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/usecase"
)

type Usecases struct {
	AuthUsecase       internal.AuthUsecase
	DiscussionUsecase internal.DiscussionUsecase
	MessageUsecase    internal.MessageUsecase
}

func newUsecases(repos *repositories) *Usecases {
	authUsecase := usecase.NewAuthUsecase(repos.userRepo)
	discussionUsecase := usecase.NewDiscussionUsecase(repos.discussionRepo)
	messageUsecase := usecase.NewMessageUsecase(repos.messageRepo)
	return &Usecases{
		AuthUsecase:       authUsecase,
		DiscussionUsecase: discussionUsecase,
		MessageUsecase:    messageUsecase,
	}
}
