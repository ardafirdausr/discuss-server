package usecase

import (
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

	msgChannel := fmt.Sprintf("%s/%v", message.ReceiverType, message.ReceiverID)
	// msgPayload, err := json.Marshal(message)
	// msgStr := string(msgPayload)
	// fmt.Println(message, msgStr)
	// if err == nil {
	// 	log.Println(err.Error())
	// 	pubsub.Publish(msgChannel, message)
	// }
	pubsub.Publish(msgChannel, []byte(message.Content))

	return message, nil
}
