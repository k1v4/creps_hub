package entity

import (
	"encoding/json"
	"time"
)

type Release struct {
	Id          int       `json:"id"`
	ReleaseDate time.Time `json:"release_date"`
	Name        string    `json:"name"`
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
