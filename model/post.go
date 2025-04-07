package model

import (
	"time"
)

type Post struct {
	ID            int       `json:"id,omitempty"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedByUser int       `json:"user_id"`
	CreatedTime   time.Time `json:"created_at"`
	UpdatedTime   time.Time `json:"udpated_at"`
}
