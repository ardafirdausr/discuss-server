package mongo

import (
	"context"
	"log"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	DB *mongo.Database
}

func NewMessageRepository(DB *mongo.Database) *MessageRepository {
	return &MessageRepository{DB: DB}
}

func (mr MessageRepository) Create(param entity.CreateMessage) (*entity.Message, error) {
	res, err := mr.DB.Collection("messages").InsertOne(context.TODO(), param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	objID := res.InsertedID.(primitive.ObjectID)
	message := &entity.Message{
		ID:        objID.Hex(),
		Content:   param.Content,
		SenderID:  param.SenderID,
		DiscussID: param.DiscussID,
		CreatedAt: param.CreatedAt,
	}

	return message, nil
}

func (mr MessageRepository) GetMessagesByDiscussionID(discussionID string) ([]*entity.Message, error) {
	return nil, nil
}
