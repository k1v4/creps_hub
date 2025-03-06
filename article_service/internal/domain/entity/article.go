package entity

import "time"

type Article struct {
	Id       int64         `json:"id"`
	Name     string        `json:"name"`
	Text     string        `json:"text"`
	Tags     []string      `json:"tags"`
	PostedAt time.Duration `json:"posted_at"`
}
