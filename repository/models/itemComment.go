package models

import (
	"github.com/go-playground/validator"
	"time"
)

type ItemComment struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId" db:"user_id"`
	ItemId    int       `json:"itemId" validate:"required" db:"item_id"`
	Comment   string    `json:"comment" validate:"required" db:"comment"`
	CreatedAt time.Time `db:"created_at"`
}

func (ic *ItemComment) Validate() error {
	validate := validator.New()

	return validate.Struct(ic)
}
