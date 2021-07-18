package ws

import (
	"errors"
	"log"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type DiscussWebSocket struct {
	app *app.App
}

func NewDiscussWebSocket(app *app.App) *DiscussWebSocket {
	dws := &DiscussWebSocket{app: app}
	return dws
}

func (dws DiscussWebSocket) authenticate(tokenizer internal.Tokenizer, strToken string) (*entity.User, error) {
	if len(strToken) < 1 {
		return nil, errors.New("token is not provided")
	}

	payload, err := tokenizer.Parse(strToken)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("invalid token")
	}

	user := &entity.User{
		ID:       payload.ID,
		Name:     payload.Name,
		Email:    payload.Email,
		ImageUrl: payload.Imageurl,
	}

	return user, nil
}
