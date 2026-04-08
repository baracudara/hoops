package dto

type Register struct {
	Name     string  `validate:"required,min=2,max=50"`
	Nickname string  `validate:"required,min=2,max=30"`
	Email    *string `validate:"omitempty,email"`
	Phone    *string `validate:"omitempty,e164"`
	GoogleID *string `validate:"omitempty"`
	Password string  `validate:"omitempty,min=8"`
}
