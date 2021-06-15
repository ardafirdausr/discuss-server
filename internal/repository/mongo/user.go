package mongo

import (
	"context"
	"log"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB mongo.Database
}

func (ur UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	res := ur.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email})
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		err := entity.ErrNotFound{
			Message: "User not found",
			Err:     res.Err(),
		}
		return nil, err
	}

	if err := res.Decode(&user); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &user, nil

}

func (ur UserRepository) Create(param entity.CreateUserParam) (*entity.User, error) {
	res, err := ur.DB.Collection("users").InsertOne(context.TODO(), param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	objID := res.InsertedID.(primitive.ObjectID)

	user := &entity.User{
		ID:       objID.Hex(),
		Email:    param.Email,
		Name:     param.Name,
		ImageUrl: param.ImageUrl,
	}
	return user, nil
}
