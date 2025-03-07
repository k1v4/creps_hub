package entity

import "encoding/json"

type User struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"-"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Username      string `json:"username"`
	AccessLevelId int    `json:"accessLevelId"`
}

func (o *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
