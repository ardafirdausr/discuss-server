package entity

import "github.com/dgrijalva/jwt-go"

type TokenPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Imageurl string `json:"imageUrl"`
}
type JWTPayload struct {
	TokenPayload
	jwt.StandardClaims
}

type GoogleAuth struct {
	TokenID  string `json:"token_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Imageurl string `json:"imageUrl,omitempty"`
}
