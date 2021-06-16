package mongo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscussionRepository struct {
	DB *mongo.Database
}

func NewDiscussionRepository(DB *mongo.Database) *DiscussionRepository {
	return &DiscussionRepository{DB: DB}
}

func (dr DiscussionRepository) GetDissionsByUserID(userID string) ([]*entity.Discussion, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return nil, err
	}

	ctx := context.TODO()
	csr, err := dr.DB.Collection("discussions").Find(ctx, bson.M{"members.id": objID})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer csr.Close(ctx)

	discussions := make([]*entity.Discussion, 0)
	for csr.Next(ctx) {
		var discussion entity.Discussion
		if err := csr.Decode(&discussion); err == nil {
			discussions = append(discussions, &discussion)
			continue
		}

		log.Println(err.Error())
	}

	return discussions, nil
}

func (dr DiscussionRepository) Create(param entity.CreateDiscussionParam) (*entity.Discussion, error) {
	param.CreatedAt = time.Now()
	param.UpdatedAt = time.Now()

	res, err := dr.DB.Collection("discussions").InsertOne(context.TODO(), param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	objID := res.InsertedID.(primitive.ObjectID)
	discussion := &entity.Discussion{
		ID:          objID.Hex(),
		Code:        param.Code,
		Name:        param.Name,
		Password:    param.Password,
		Description: param.Description,
		PhotoUrl:    param.PhotoUrl,
		CreatedAt:   param.CreatedAt,
		UpdatedAt:   param.UpdatedAt,
	}

	return discussion, nil
}

func (dr DiscussionRepository) Update(ID string, param entity.CreateDiscussionParam) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return false, err
	}

	param.UpdatedAt = time.Now()

	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), objID, param)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	if res.MatchedCount < 1 {
		return false, errors.New("failed to update data")
	}

	return true, nil
}

func (dr DiscussionRepository) DeleteByID(discussionID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(discussionID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return false, err
	}

	ctx := context.TODO()
	res, err := dr.DB.Collection("discussions").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	if res.DeletedCount < 1 {
		return false, errors.New("failed to delete data")
	}

	return true, nil
}
