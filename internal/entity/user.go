package entity

type User struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	ImageUrl    string        `json:"image_url"`
	Discussions []*Discussion `json:"discuses"`
}

type CreateUserParam struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"image_url"`
}
