package entity

import (
	"encoding/json"
	"time"
)

type Article struct {
	ID              int       `json:"id"`
	AuthorID        int       `json:"author_id"`
	PublicationDate time.Time `json:"publication_date"`
	Name            string    `json:"name"`
	Text            string    `json:"text"`
}

func (o *Article) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Article) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
