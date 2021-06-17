package entity

type User struct {
	ID          interface{}   `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Email       string        `json:"email" bson:"email"`
	ImageUrl    string        `json:"image_url" bson:"imageUrl"`
	Discussions []*Discussion `json:"discussions" bson:"discussions"`
}

type CreateUserParam struct {
	Name        string        `json:"name" bson:"name"`
	Email       string        `json:"email" bson:"email"`
	ImageUrl    string        `json:"image_url" bson:"imageUrl"`
	Discussions []*Discussion `bson:"discussions"`
}
