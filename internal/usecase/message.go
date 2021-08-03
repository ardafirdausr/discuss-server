package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type MessageUsecase struct {
	messageRepo internal.MessageRepository
}

func NewMessageUsecase(messageRepo internal.MessageRepository) *MessageUsecase {
	return &MessageUsecase{messageRepo: messageRepo}
}

func (muc MessageUsecase) GetMessagesByDiscussionID(discussionID interface{}, size, page int) ([]*entity.Message, error) {
	if size < 20 {
		size = 20
	}

	if page < 1 {
		size = 1
	}

	messages, err := muc.messageRepo.GetPaginatedMessagesByDiscussionID(discussionID, size, page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return messages, nil
}

func (muc MessageUsecase) SendMessage(pubsub internal.PubSub, param entity.CreateMessage) (*entity.Message, error) {
	param.CreatedAt = time.Now()
	message, err := muc.messageRepo.Create(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	msgChannel := fmt.Sprintf("%s/%v", message.ReceiverType, message.ReceiverID)
	msgPayload, err := json.Marshal(&message)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	pubsub.Publish(msgChannel, string(msgPayload))
	return message, nil
}
