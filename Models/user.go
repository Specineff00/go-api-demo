package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`     // Make sure formatting of json correct
	Name      string    `json:"name" db:"name"` // TODO: struct tags for formatting json
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
