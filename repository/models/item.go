package models

import (
	"github.com/go-playground/validator"
)

type Item struct {
	Id           int      `json:"id"`
	UserId       int      `db:"user_id"`
	Title        string   `json:"title" validate:"required"`
	InitialPrice float64  `json:"initialPrice" validate:"min=0" db:"initial_price"`
	SoldPrice    *float64 `json:"soldPrice" db:"sold_price"`
	Description  *string  `json:"description"`
}

func (i *Item) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}
