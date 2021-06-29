package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type MessageUsecase struct {
	messageRepo internal.MessageRepository
}

func NewMessageUsecase(messageRepo internal.MessageRepository) *MessageUsecase {
	return &MessageUsecase{messageRepo: messageRepo}
}

func (muc MessageUsecase) SendMessage(pubsub internal.PubSub, param entity.CreateMessage) (*entity.Message, error) {
	message, err := muc.messageRepo.Create(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	msgEvt := entity.NewMessageSent(*message)
	msgChannel := fmt.Sprintf("%s/%v", message.ReceiverType, message.ReceiverID)
	msgPayload, err := json.Marshal(msgEvt)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	pubsub.Publish(msgChannel, msgPayload)
	return message, nil
}
