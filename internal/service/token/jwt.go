package token

import (
	"errors"
	"log"
	"time"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

type JWTTokenizer struct {
	secretKey string
}

func NewJWTTokenizer(secretKey string) JWTTokenizer {
	return JWTTokenizer{secretKey: secretKey}
}

func (JWTT JWTTokenizer) Generate(payload entity.TokenPayload) (string, error) {
	jwtPayload := entity.JWTPayload{}
	jwtPayload.TokenPayload = payload
	jwtPayload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)
	jwtToken, err := token.SignedString([]byte(JWTT.secretKey))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return jwtToken, nil
}

func (JWTT JWTTokenizer) Parse(tokenString string) (*entity.TokenPayload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(JWTT.secretKey), nil
	}

	claims := entity.JWTPayload{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !token.Valid {
		err = errors.New("invalid token")
		log.Println(err.Error())
		return nil, err
	}

	return &claims.TokenPayload, nil
}
