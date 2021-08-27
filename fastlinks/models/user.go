package models

import "time"

type User struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string

	// TODO: perms fields
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
