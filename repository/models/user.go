package models

import (
	"github.com/go-playground/validator"
	"time"
)

type User struct {
	Id           int       `json:"id"`
	FirstName    string    `json:"firstName" validate:"required,min=1" db:"first_name"`
	LastName     string    `json:"lastName" validate:"required,min=1" db:"last_name"`
	Email        string    `json:"email" validate:"required,email,min=1"`
	PasswordHash string    `json:"password" db:"password_hash"`
	LastLoginUtc time.Time `db:"last_login_utc"`
	UserTypeId   int32     `json:"userTypeId" validate:"required" db:"user_type_id"`
}

func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}
