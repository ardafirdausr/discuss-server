package mongo

import (
	"context"
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageModel struct {
	ID           primitive.ObjectID         `bson:"_id,omitempty"`
	ContentType  entity.MessageContentType  `bson:"contentType,omitempty"`
	Content      string                     `bson:"content,omitempty"`
	ReceiverType entity.MessageReceiverType `bson:"receiverType,omitempty"`
	ReceiverID   primitive.ObjectID         `bson:"receiverId,omitempty"`
	Sender       userModel                  `bson:"sender,omitempty"`
	CreatedAt    time.Time                  `bson:"createdAt,omitempty"`
}

func (mm *messageModel) toMessage() *entity.Message {
	message := &entity.Message{
		ID:           mm.ID.Hex(),
		ContentType:  mm.ContentType,
		Content:      mm.Content,
		ReceiverType: mm.ReceiverType,
		ReceiverID:   mm.ReceiverID.Hex(),
		Sender:       *mm.Sender.ToUser(),
		CreatedAt:    mm.CreatedAt,
	}

	return message
}

type MessageRepository struct {
	DB *mongo.Database
}

func NewMessageRepository(DB *mongo.Database) *MessageRepository {
	return &MessageRepository{DB: DB}
}

func (mr MessageRepository) Create(param entity.CreateMessage) (*entity.Message, error) {

	senderObjID, err := toObjectID(param.Sender.ID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	userModel := userModel{
		ID:       senderObjID,
		Name:     param.Sender.Name,
		Email:    param.Sender.Email,
		ImageUrl: param.Sender.ImageUrl,
	}

	receiverObjID, err := toObjectID(param.ReceiverID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	messageModel := &messageModel{
		ContentType:  param.ContentType,
		Content:      param.Content,
		ReceiverType: param.ReceiverType,
		ReceiverID:   receiverObjID,
		Sender:       userModel,
		CreatedAt:    param.CreatedAt,
	}

	res, err := mr.DB.Collection("messages").InsertOne(context.TODO(), messageModel)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	messageModel.ID = res.InsertedID.(primitive.ObjectID)

	return messageModel.toMessage(), nil
}

func (mr MessageRepository) GetMessagesByDiscussionID(discussionID interface{}) ([]*entity.Message, error) {
	return nil, nil
}
