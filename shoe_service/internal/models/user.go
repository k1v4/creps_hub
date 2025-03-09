package models

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	UserName string `json:"username"`
	Shoes    []Shoe `json:"shoes"`
}
