package models

import (
	"github.com/go-playground/validator"
)

type Item struct {
	Id           int      `json:"id"`
	UserId       int      `json:"userId" db:"user_id"`
	Title        string   `json:"title" validate:"required"`
	InitialPrice float64  `json:"initialPrice" validate:"min=0" db:"initial_price"`
	SoldPrice    *float64 `json:"soldPrice" db:"sold_price"`
	Description  *string  `json:"description"`
}

func (i *Item) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}

//func (i *Item) MarshalJSON() ([]byte, error) {
//	type Alias Item
//	return json.Marshal(&struct {
//		SoldPrice float64 `json:"soldPrice"`
//		*Alias
//	}{
//		SoldPrice: func() float64 {
//			if i.SoldPrice.Valid {
//				return i.SoldPrice.Float64
//			}
//			return 0.0
//		}(),
//		Alias: (*Alias)(i),
//	})
//}
