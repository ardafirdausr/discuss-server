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
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name,omitempty"`
	Email       string               `bson:"email,omitempty"`
	ImageUrl    string               `bson:"imageUrl,omitempty"`
	Discussions []primitive.ObjectID `bson:"Discussions,omitempty"`
}

func (um *userModel) ToUser() *entity.User {
	return &entity.User{
		ID:       um.ID,
		Name:     um.Name,
		Email:    um.Email,
		ImageUrl: um.ImageUrl,
	}
}

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(DB *mongo.Database) *UserRepository {
	return &UserRepository{DB: DB}
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
