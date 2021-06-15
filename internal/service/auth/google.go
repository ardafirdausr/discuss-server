package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"google.golang.org/api/idtoken"
)

type GoogleSSOAuthenticator struct {
	clientID string
}

func NewGoogleSSOAuthenticator(clientID string) GoogleSSOAuthenticator {
	return GoogleSSOAuthenticator{clientID: clientID}
}

func (auth GoogleSSOAuthenticator) getTokenPayload(token string) (*idtoken.Payload, error) {
	clientOption := idtoken.WithHTTPClient(&http.Client{})
	tokenValidator, err := idtoken.NewValidator(context.Background(), clientOption)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	payload, err := tokenValidator.Validate(context.Background(), token, auth.clientID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return payload, nil
}

func (auth GoogleSSOAuthenticator) Authenticate(token string) (*entity.User, error) {
	payload, err := auth.getTokenPayload(token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	user := &entity.User{
		Name:     payload.Claims["name"].(string),
		Email:    payload.Claims["email"].(string),
		ImageUrl: payload.Claims["picture"].(string),
	}

	return user, nil
}
