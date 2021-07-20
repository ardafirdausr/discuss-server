package mongo

import (
	"context"
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (mr MessageRepository) GetPaginatedMessagesByDiscussionID(discussionID interface{}, size int, page int) ([]*entity.Message, error) {
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	ctx := context.TODO()
	filter := bson.M{"receiverType": "message.receiver.discussion", "receiverId": discussionObjID}
	skip := int64(size * (page - 1))
	limit := int64(size)
	options := &options.FindOptions{}
	options.SetLimit(limit)
	options.SetSkip(skip)
	options.SetSort(bson.M{"createdAt": -1})
	csr, err := mr.DB.Collection("messages").Find(ctx, filter, options)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer csr.Close(ctx)

	var messages []*entity.Message
	for csr.Next(ctx) {
		var messageModel messageModel
		if err := csr.Decode(&messageModel); err != nil {
			log.Println(err.Error())
			continue
		}

		messages = append(messages, messageModel.toMessage())
	}
	log.Println(messages)

	return messages, nil
}
