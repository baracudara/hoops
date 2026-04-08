package storage

import "errors"

var (

	// sql
	ErrUserNotFound = errors.New("User not found")
	ErrUserExists = errors.New("User alredy exists")
	ErrAppNotFound = errors.New("App not found")

	// redis

	ErrTokenNotFound  = errors.New("Token not found")

	
)