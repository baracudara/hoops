package models

import "time"

type User struct {
	ID       string
	Name     string
	Nickname string
	Email    *string
	Phone    *string
	GoogleID *string
	PassHash []byte
	Role Role
	TrustRating int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}



