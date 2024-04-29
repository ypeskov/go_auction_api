package models

import "github.com/go-playground/validator"

type Item struct {
	Id           int    `json:"id"`
	Title        string `json:"title" validate:"required"`
	InitialPrice int    `json:"initialPrice" validate:"min=0"`
	SoldPrice    int    `json:"soldPrice"`
	Description  string `json:"description"`
}

func (i *Item) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}
