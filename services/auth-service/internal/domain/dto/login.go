package dto

type Login struct {
	Email    *string `validate:"omitempty,email"`
	Phone    *string `validate:"omitempty,e164"`
	GoogleID *string `validate:"omitempty"`
	Password string  `validate:"omitempty,min=8"`
}
