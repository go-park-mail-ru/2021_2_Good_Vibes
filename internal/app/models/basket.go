package models

type BasketProduct struct {
	UserId    int `json:"user_id,omitempty"`
	ProductId int `json:"product_id" validate:"required"`
	Number    int `json:"number,omitempty"`
}

type BasketStorage struct {
	UserId int `json:"user_id"`
}
