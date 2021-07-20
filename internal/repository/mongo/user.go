package mongo

import (
	"context"
	"log"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userModel struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name,omitempty"`
	Email         string               `bson:"email,omitempty"`
	ImageUrl      string               `bson:"imageUrl,omitempty"`
	DiscussionIDs []primitive.ObjectID `bson:"discussionIds,omitempty"`

	Discussions []*discussionModel `bson:"Discussions,omitempty"`
}

func (um *userModel) ToUser() *entity.User {
	user := &entity.User{
		ID:       um.ID.Hex(),
		Name:     um.Name,
		Email:    um.Email,
		ImageUrl: um.ImageUrl,
	}

	if len(um.Discussions) > 0 {
		var discussions []*entity.Discussion
		for _, discussionModel := range um.Discussions {
			discussions = append(discussions, discussionModel.toDiscussion())
		}

		user.Discussions = discussions
	}

	return user
}

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(DB *mongo.Database) *UserRepository {
	return &UserRepository{DB: DB}
}

func (ur UserRepository) loadDiscussions(um *userModel) error {
	if len(um.DiscussionIDs) < 1 {
		return nil
	}

	ctx := context.TODO()
	csr, err := ur.DB.Collection("discussions").Find(ctx, bson.M{"_id": bson.M{"$in": um.DiscussionIDs}})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer csr.Close(ctx)

	for csr.Next(ctx) {
		var discussionModel discussionModel
		if err := csr.Decode(&discussionModel); err != nil {
			log.Println(err.Error())
			continue
		}

		um.Discussions = append(um.Discussions, &discussionModel)
	}

	return nil
}

func (ur UserRepository) GetUserByID(userID interface{}) (*entity.User, error) {
	userObjID, err := toObjectID(userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var userModel userModel
	res := ur.DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": userObjID})
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		err := entity.ErrNotFound{
			Message: "User not found",
			Err:     res.Err(),
		}
		return nil, err
	}

	if err := res.Decode(&userModel); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = ur.loadDiscussions(&userModel)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return userModel.ToUser(), nil
}

func (ur UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var userModel userModel
	res := ur.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email})
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		err := entity.ErrNotFound{
			Message: "User not found",
			Err:     res.Err(),
		}
		return nil, err
	}

	if err := res.Decode(&userModel); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err := ur.loadDiscussions(&userModel)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return userModel.ToUser(), nil
}

func (ur UserRepository) Create(param entity.CreateUserParam) (*entity.User, error) {
	userModel := userModel{
		Email:    param.Email,
		Name:     param.Name,
		ImageUrl: param.ImageUrl,
	}

	res, err := ur.DB.Collection("users").InsertOne(context.TODO(), userModel)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	userModel.ID = res.InsertedID.(primitive.ObjectID)
	return userModel.ToUser(), nil
}
