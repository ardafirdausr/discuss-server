package usecase

import (
	"errors"
	"log"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type AuthUsecase struct {
	userRepo internal.UserRepository
}

func NewAuthUsecase(userRepo internal.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (service AuthUsecase) SSO(token string, authenticator internal.SSOAuthenticator) (*entity.User, error) {
	reqUser, err := authenticator.Authenticate(token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	user, err := service.userRepo.GetUserByEmail(reqUser.Email)
	if _, ok := err.(entity.ErrNotFound); ok {
		param := entity.CreateUserParam{
			Email:    reqUser.Email,
			Name:     reqUser.Name,
			ImageUrl: reqUser.ImageUrl,
		}
		user, err = service.userRepo.Create(param)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return user, nil
}

func (service AuthUsecase) GetUserFromToken(token string, tokenizer internal.Tokenizer) (*entity.User, error) {
	if len(token) < 1 {
		return nil, errors.New("token is not provided")
	}

	payload, err := tokenizer.Parse(token)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("invalid token")
	}

	user, err := service.userRepo.GetUserByID(payload.ID)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (service AuthUsecase) GenerateAuthToken(user entity.User, tokenizer internal.Tokenizer) (string, error) {
	tokenPayload := entity.TokenPayload{}
	tokenPayload.ID = user.ID
	tokenPayload.Name = user.Name
	tokenPayload.Email = user.Email
	tokenPayload.Imageurl = user.ImageUrl
	token, err := tokenizer.Generate(tokenPayload)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return token, nil
}
