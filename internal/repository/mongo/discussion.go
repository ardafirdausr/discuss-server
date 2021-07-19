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

type discussionModel struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Code        string               `bson:"code,omitempty"`
	Name        string               `bson:"name,omitempty"`
	Description string               `bson:"description,omitempty"`
	Password    *string              `bson:"password,omitempty"`
	PhotoUrl    *string              `bson:"photoUrl,omitempty"`
	CreatorID   primitive.ObjectID   `bson:"creatorId,omitempty"`
	MemberIDs   []primitive.ObjectID `bson:"memberIds,omitempty"`
	CreatedAt   time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt   time.Time            `bson:"updatedAt,omitempty"`

	Members  []*userModel    `bson:"members,omitempty"`
	Messages []*messageModel `bson:"messages,omitempty"`
}

func (dm *discussionModel) toDiscussion() *entity.Discussion {
	discussion := &entity.Discussion{
		ID:          dm.ID.Hex(),
		Code:        dm.Code,
		Name:        dm.Name,
		Description: dm.Description,
		Password:    dm.Password,
		PhotoUrl:    dm.PhotoUrl,
		CreatorID:   dm.CreatorID,
		CreatedAt:   dm.CreatedAt,
		UpdatedAt:   dm.UpdatedAt,
	}

	if len(dm.Members) > 0 {
		var users []*entity.User
		for _, userModel := range dm.Members {
			users = append(users, userModel.ToUser())
		}

		discussion.Members = users
	}

	return discussion
}

type DiscussionRepository struct {
	DB *mongo.Database
}

func NewDiscussionRepository(DB *mongo.Database) *DiscussionRepository {
	return &DiscussionRepository{DB: DB}
}

func (dr DiscussionRepository) loadMembers(model *discussionModel) error {
	if len(model.MemberIDs) < 1 {
		return nil
	}

	ctx := context.TODO()
	csr, err := dr.DB.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": model.MemberIDs}})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer csr.Close(ctx)

	model.Members = make([]*userModel, 0)
	for csr.Next(ctx) {
		var userModel userModel
		if err := csr.Decode(&userModel); err == nil {
			model.Members = append(model.Members, &userModel)
			continue
		}

		log.Println(err.Error())
	}

	return nil
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

	var model discussionModel
	if err := res.Decode(&model); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if err := dr.loadMembers(&model); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return model.toDiscussion(), nil
}

func (dr DiscussionRepository) GetDiscussionByCode(code string) (*entity.Discussion, error) {
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

	var model discussionModel
	if err := res.Decode(&model); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if err := dr.loadMembers(&model); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return model.toDiscussion(), nil
}

func (dr DiscussionRepository) GetDiscussionsByUserID(userID interface{}) ([]*entity.Discussion, error) {
	objID, err := toObjectID(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	matchStage := bson.M{"$match": bson.M{"memberIds": objID}}
	lookupStage := bson.M{"$lookup": bson.M{"from": "users", "localField": "memberIds", "foreignField": "_id", "as": "members"}}
	ctx := context.TODO()
	csr, err := dr.DB.Collection("discussions").Aggregate(context.TODO(), []bson.M{matchStage, lookupStage})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer csr.Close(ctx)

	var discussions []*entity.Discussion
	for csr.Next(ctx) {
		var discussionModel discussionModel
		if err := csr.Decode(&discussionModel); err == nil {
			discussion := discussionModel.toDiscussion()
			discussions = append(discussions, discussion)
			continue
		}

		log.Println(err.Error())
	}

	return discussions, nil
}

func (dr DiscussionRepository) Create(param entity.CreateDiscussionParam) (*entity.Discussion, error) {
	creatorObjID, err := toObjectID(param.CreatorID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	model := discussionModel{
		Code:        param.Code,
		Name:        param.Name,
		Description: param.Description,
		Password:    param.Password,
		CreatorID:   creatorObjID,
		MemberIDs:   []primitive.ObjectID{creatorObjID},
		CreatedAt:   param.CreatedAt,
		UpdatedAt:   param.UpdatedAt,
	}

	res, err := dr.DB.Collection("discussions").InsertOne(context.TODO(), model)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	model.ID = res.InsertedID.(primitive.ObjectID)
	if err := dr.loadMembers(&model); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	update := bson.M{"$addToSet": bson.M{"discussionIds": model.ID}}
	_, err = dr.DB.Collection("users").UpdateByID(context.TODO(), creatorObjID, update)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return model.toDiscussion(), nil
}

func (dr DiscussionRepository) AddMember(discussionID interface{}, userID interface{}) error {
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
	update := bson.M{"$addToSet": bson.M{"memberIds": userObjID}}
	res, err := dr.DB.Collection("discussions").UpdateByID(ctx, discussionObjID, update)
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

	update = bson.M{"$addToSet": bson.M{"discussionIds": discussionObjID}}
	res, err = dr.DB.Collection("users").UpdateByID(context.TODO(), userObjID, update)
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
	res, err := dr.DB.Collection("discussions").UpdateByID(ctx, discussionObjID, bson.M{"$pull": bson.M{"memberIds": userObjID}})
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
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	model := discussionModel{
		Code:        param.Code,
		Name:        param.Name,
		Description: param.Description,
		UpdatedAt:   param.UpdatedAt,
	}
	update := bson.M{"$set": model}
	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), discussionObjID, update)
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
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	data := bson.M{"password": password}
	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), discussionObjID, bson.M{"$set": data})
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
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	data := bson.M{"photoUrl": url}
	res, err := dr.DB.Collection("discussions").UpdateByID(context.TODO(), discussionObjID, bson.M{"$set": data})
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
	discussionObjID, err := toObjectID(discussionID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	ctx := context.TODO()
	res, err := dr.DB.Collection("discussions").DeleteOne(ctx, bson.M{"_id": discussionObjID})
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
