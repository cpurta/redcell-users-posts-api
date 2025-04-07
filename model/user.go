package model

import (
	"time"
)

type User struct {
	ID          int       `json:"id,omitempty"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	TimeCreated time.Time `json:"created_at,omitempty"`
	TimeUpdated time.Time `json:"updated_at,omitempty"`
}
