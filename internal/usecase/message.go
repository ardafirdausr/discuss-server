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

func (muc MessageUsecase) SendMessage(pubsub internal.PubSub, param entity.CreateMessage) (*entity.Message, error) {
	param.CreatedAt = time.Now()
	message, err := muc.messageRepo.Create(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	msgChannel := fmt.Sprintf("%s/%v", message.ReceiverType, message.ReceiverID)
	fmt.Println("WWWWWWWWWWWWWWWWWWWWWWWWW")
	fmt.Println(msgChannel)
	fmt.Println("WWWWWWWWWWWWWWWWWWWWWWWWW")
	msgPayload, err := json.Marshal(&message)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	pubsub.Publish(msgChannel, msgPayload)
	return message, nil
}
