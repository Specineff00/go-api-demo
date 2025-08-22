package models

type User struct {
	ID   int    `json:"id"`   // Make sure formatting of json correct
	Name string `json:"name"` // TODO: struct tags for formatting json
}
