package entity

import "time"

type Discussion struct {
	ID          interface{} `json:"id"`
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Password    *string     `json:"-"`
	PhotoUrl    *string     `json:"image_url"`
	CreatorID   interface{} `json:"creator_id"`
	Members     []*User     `json:"members"`
	Messages    []*Message  `json:"messages"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CreateDiscussionParam struct {
	Code        string      `json:"code" validate:"required,min=4"`
	Name        string      `json:"name" validate:"required,min=4"`
	Description string      `json:"description"`
	Password    *string     `json:"password"`
	CreatorID   interface{} `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type JoinDiscussionParam struct {
	Code     string      `json:"code" validate:"required"`
	Password *string     `json:"password"`
	UserID   interface{} `validate:"required"`
}

type UpdateDiscussionParam struct {
	Code        string `json:"code,omitempty" validate:"min=4"`
	Name        string `json:"name,omitempty" validate:"min=4"`
	Description string `json:"description,omitempty"`
	UpdatedAt   time.Time
}

type UpdateDiscussionPassword struct {
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
