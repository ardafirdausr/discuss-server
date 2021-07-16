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

func (dr DiscussionRepository) GetDiscussionsByID(discussionID interface{}) (*entity.Discussion, error) {
	objID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
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
	objID, err := toObjectID(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	matchStage := bson.M{"$match": bson.M{"members": objID}}
	lookupStage := bson.M{"$lookup": bson.M{"from": "users", "localField": "members", "foreignField": "_id", "as": "members"}}
	ctx := context.TODO()
	csr, err := dr.DB.Collection("discussions").Aggregate(context.TODO(), []bson.M{matchStage, lookupStage})
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
	objID, err := toObjectID(param.CreatorID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	param.CreatorID = objID
	var objMemberIDs []interface{}
	for _, memberID := range param.Members {
		memberIDStr := memberID.(string)
		objmemberID, err := primitive.ObjectIDFromHex(memberIDStr)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		objMemberIDs = append(objMemberIDs, objmemberID)
	}
	param.Members = objMemberIDs

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

func (dr DiscussionRepository) AddMember(discussionID, userID interface{}) error {
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	userObjID, err := toObjectID(userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ctx := context.TODO()
	res, err := dr.DB.Collection("discussions").UpdateByID(ctx, discussionObjID, bson.M{"$addToSet": bson.M{"members": userObjID}})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.MatchedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) RemoveMember(discussionID, userID interface{}) error {
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	userObjID, err := toObjectID(userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ctx := context.TODO()
	res, err := dr.DB.Collection("discussions").UpdateByID(ctx, discussionObjID, bson.M{"$pull": bson.M{"members": userObjID}})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.MatchedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) UpdateByID(discussionID interface{}, param entity.UpdateDiscussionParam) error {
	objID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
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
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) UpdatePasswordByID(discussionID interface{}, password string) error {
	objID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	data := bson.M{"password": password}
	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), objID, bson.M{"$set": data})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.MatchedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) UpdatePhotoByID(discussionID interface{}, url string) error {
	objID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	data := bson.M{"photoUrl": url}
	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), objID, bson.M{"$set": data})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.MatchedCount < 1 {
		err = entity.ErrNotFound{
			Message: "Discussion not found",
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}

func (dr DiscussionRepository) DeleteByID(discussionID interface{}) error {
	objID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
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
			Err:     errors.New("document not found"),
		}
		return err
	}

	return nil
}
