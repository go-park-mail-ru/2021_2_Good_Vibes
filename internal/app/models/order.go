package models

type OrderProducts struct {
	OrderId int `json:"order_id,omitempty"`
	ProductId int `json:"product_id"`
	Number int `json:"number"`
}

type Order struct {
	OrderId int `json:"order_id,omitempty"`
	UserId    int  `json:"user_id"`
	Date     string  `json:"date"`
	Address string `json:"address"`
	Cost   float32 `json:"cost"`
	Status string  `json:"status"`
	Products []OrderProducts `json:"products"`
}
