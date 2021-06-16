package entity

import "time"

type Discussion struct {
	ID          string    `json:"id" bson:"id"`
	Code        string    `json:"string" bson:"string"`
	Name        string    `json:"name" bson:"name"`
	Password    *string   `json:"-" bson:"password"`
	Description string    `json:"description" bson:"description"`
	PhotoUrl    string    `json:"photo_url" bson:"photoUrl"`
	CreatorID   string    `json:"creator" bson:"creator"`
	Member      []*User   `json:"members" bson:"members"`
	CreatedAt   time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updatedAt"`
}

type CreateDiscussionParam struct {
	Code        string    `json:"code" bson:"code"`
	Name        string    `json:"name" bson:"name"`
	Password    *string   `json:"password,omitempty" bson:"password"`
	Description string    `json:"description" bson:"description"`
	PhotoUrl    string    `json:"photo_url" bson:"photoUrl"`
	CreatorID   string    `json:"creator" bson:"creator"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

type UpdateDiscussionParam struct {
	Code        string    `json:"code" bson:"code"`
	Name        string    `json:"name" bson:"name"`
	Password    *string   `json:"password,omitempty" bson:"password"`
	Description string    `json:"description" bson:"description"`
	PhotoUrl    string    `json:"photo_url" bson:"photoUrl"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
