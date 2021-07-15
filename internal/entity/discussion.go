package entity

import "time"

type Discussion struct {
	ID          interface{} `json:"id" bson:"_id"`
	Code        string      `json:"code" bson:"code"`
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Password    *string     `json:"-" bson:"password"`
	PhotoUrl    *string     `json:"image_url" bson:"photoUrl"`
	CreatorID   interface{} `json:"creator_id" bson:"creatorId"`
	Members     []*User     `json:"members" bson:"members"`
	CreatedAt   time.Time   `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updatedAt"`
}

type CreateDiscussionParam struct {
	Code        string        `json:"code" validate:"required,min=4" bson:"code"`
	Name        string        `json:"name" validate:"required,min=4" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Password    *string       `json:"-" bson:"password"`
	CreatorID   interface{}   `bson:"creatorId" validate:"required"`
	Members     []interface{} `bson:"members"`
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
}

type UpdateDiscussionParam struct {
	Code        string    `json:"code,omitempty" validate:"min=4" bson:"code"`
	Name        string    `json:"name,omitempty" validate:"min=4" bson:"name"`
	Description string    `json:"description,omitempty" bson:"description"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type UpdateDiscussionPassword struct {
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
