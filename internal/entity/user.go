package entity

type User struct {
	ID          interface{}  `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	ImageUrl    string       `json:"image_url"`
	Discussions []Discussion `json:"Discussions,omitempty"`
}

type CreateUserParam struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"image_url"`
}
