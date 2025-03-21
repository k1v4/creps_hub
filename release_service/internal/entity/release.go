package entity

import (
	"encoding/json"
	"time"
)

type Release struct {
	Id          int       `json:"id"`
	ReleaseDate time.Time `json:"release_date"`
	Name        string    `json:"name"`
	ImageUrl    string    `json:"image_url"`
}

func (o *Release) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Release) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Password      []byte `json:"-"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Username      string `json:"username"`
	AccessLevelId int    `json:"-"`
}
