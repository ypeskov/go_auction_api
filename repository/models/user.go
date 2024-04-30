package models

import (
	"github.com/go-playground/validator"
	"time"
)

type User struct {
	Id           int       `json:"id"`
	FirsName     string    `json:"firstName" validate:"required,min=1" db:"first_name"`
	LastName     string    `json:"lastName" validate:"required,min=1" db:"last_name"`
	Email        string    `json:"email" validate:"required,email,min=1"`
	LastLoginUtc time.Time `db:"last_login_utc"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}