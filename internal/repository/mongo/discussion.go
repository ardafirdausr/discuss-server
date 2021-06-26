package mongo

import (
	"context"
	"errors"
	"log"

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

func (dr DiscussionRepository) GetDiscussionsByID(ID interface{}) (*entity.Discussion, error) {
	strID := ID.(string)
	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		log.Println(err.Error())
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     err,
		}
		return nil, err
	}

	ctx := context.TODO()
	res := dr.DB.Collection("discussions").FindOne(ctx, bson.M{"_id": objID})
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		err := entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     res.Err(),
		}
		return nil, err
	}

	var discussion entity.Discussion
	if err := res.Decode(&discussion); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &discussion, nil
}

func (dr DiscussionRepository) GetDiscussionsByCode(code string) (*entity.Discussion, error) {
	ctx := context.TODO()
	res := dr.DB.Collection("discussions").FindOne(ctx, bson.M{"code": code})
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		err := entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     res.Err(),
		}
		return nil, err
	}

	var discussion entity.Discussion
	if err := res.Decode(&discussion); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &discussion, nil
}

func (dr DiscussionRepository) GetDiscussionsByUserID(userID interface{}) ([]*entity.Discussion, error) {
	strUserID, ok := userID.(string)
	if !ok {
		err := entity.ErrNotFound{Message: "Invalid ID"}
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(strUserID)
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
	res, err := dr.DB.Collection("discussions").InsertOne(context.TODO(), param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	discussion := &entity.Discussion{
		ID:          res.InsertedID,
		Code:        param.Code,
		Name:        param.Name,
		Description: param.Description,
		PhotoUrl:    nil,
		CreatorID:   param.CreatorID,
		CreatedAt:   param.CreatedAt,
		UpdatedAt:   param.UpdatedAt,
	}

	return discussion, nil
}

func (dr DiscussionRepository) UpdateByID(ID interface{}, param entity.UpdateDiscussionParam) error {
	strID, ok := ID.(string)
	if !ok {
		err := entity.ErrNotFound{Message: "Invalid ID"}
		return err
	}

	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return err
	}

	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), objID, bson.M{"$set": param})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.MatchedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("Document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) DeleteByID(discussionID interface{}) error {
	strDiscussionID, ok := discussionID.(string)
	if !ok {
		err := entity.ErrNotFound{Message: "Invalid ID"}
		return err
	}

	objID, err := primitive.ObjectIDFromHex(strDiscussionID)
	if err != nil {
		err = entity.ErrNotFound{
			Message: "Failed get data using the corresponding ID",
			Err:     err,
		}
		return err
	}

	ctx := context.TODO()
	res, err := dr.DB.Collection("discussions").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.DeletedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("Document not found"),
		}
		return err
	}

	return nil
}
