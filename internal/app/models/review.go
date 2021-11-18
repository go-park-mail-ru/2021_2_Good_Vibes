package models


type Review struct {
	UserId int `json:"-"`
	ProductId int `json:"product_id" validate:"required"`
	Rating int `json:"rating" validate:"required"`
	Text string `json:"text" validate:"required"`
}
