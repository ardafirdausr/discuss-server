package app

import (
	"github.com/ardafirdausr/discuss-server/internal"
	mongoRepo "github.com/ardafirdausr/discuss-server/internal/repository/mongo"
)

type repositories struct {
	userRepo       internal.UserRepository
	discussionRepo internal.DiscussionRepository
	messageRepo    internal.MessageRepository
}

func newRepositories(drivers *drivers) *repositories {
	userRepo := mongoRepo.NewUserRepository(drivers.Mongo)
	discussionRepo := mongoRepo.NewDiscussionRepository(drivers.Mongo)
	messageRepo := mongoRepo.NewMessageRepository(drivers.Mongo)

	return &repositories{
		userRepo:       userRepo,
		discussionRepo: discussionRepo,
		messageRepo:    messageRepo,
	}
}
