package models

type Shoe struct {
	Id       int64  `json:"shoe_id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	UserId   int64  `json:"user_id"`
}
