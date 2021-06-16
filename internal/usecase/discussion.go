package usecase

import (
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type DiscussionUsecase struct {
	discussionRepo internal.DiscussionRepository
	messageRepo    internal.MessageRepository
}

func NewDiscussionUsecase(
	discussionRepo internal.DiscussionRepository,
	messageRepo internal.MessageRepository) *DiscussionUsecase {
	return &DiscussionUsecase{
		discussionRepo: discussionRepo,
		messageRepo:    messageRepo,
	}
}

func (du DiscussionUsecase) GetAllUserDiscussions(userID string) ([]*entity.Discussion, error) {
	discussions, err := du.discussionRepo.GetDissionsByUserID(userID)
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}

	return discussions, nil
}

func (du DiscussionUsecase) GetDiscussionMessages(dicussionID string) ([]*entity.Message, error) {
	messages, err := du.messageRepo.GetMessagesByDiscussionID(dicussionID)
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}

	return messages, nil
}

func (du DiscussionUsecase) Create(param entity.CreateDiscussionParam) (*entity.Discussion, error) {
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()
	discussion, err := du.discussionRepo.Create(param)
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}

	return discussion, nil
}

func (du DiscussionUsecase) SendMessage(dicussionID string, param entity.CreateMessage) (*entity.Message, error) {
	param.DiscussID = dicussionID
	param.CreatedAt = time.Now()
	message, err := du.messageRepo.Create(param)
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}

	return message, nil
}

func (du DiscussionUsecase) Update(dicussionID string, param entity.UpdateDiscussionParam) (bool, error) {
	param.UpdatedAt = time.Now()
	isUpdated, err := du.discussionRepo.UpdateByID(dicussionID, param)
	if err != nil {
		log.Panicln(err.Error())
		return false, err
	}

	return isUpdated, nil
}

func (du DiscussionUsecase) Delete(dicussionID string) (bool, error) {
	isDeleted, err := du.discussionRepo.DeleteByID(dicussionID)
	if err != nil {
		log.Panicln(err.Error())
		return false, err
	}

	return isDeleted, nil
}
