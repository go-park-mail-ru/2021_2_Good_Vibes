package models

type OrderProducts struct {
	OrderId   int `json:"order_id,omitempty"`
	ProductId int `json:"product_id" validate:"required"`
	Number    int `json:"number" validate:"required"`
}

type Address struct {
	Country string `json:"country" validate:"required"`
	Region  string `json:"region" validate:"required"`
	City    string `json:"city" validate:"required"`
	Street  string `json:"street" validate:"required"`
	House   string `json:"house" validate:"required"`
	Flat    string `json:"flat" validate:"required"`
	Index   string `json:"index" validate:"required"`
}

type Order struct {
	OrderId  int             `json:"order_id,omitempty"`
	UserId   int             `json:"user_id,omitempty"`
	Date     string          `json:"date,omitempty"`
	Address  Address         `json:"address" validate:"required"`
	Cost     float64         `json:"cost,omitempty"`
	Status   string          `json:"status,omitempty"`
	Products []OrderProducts `json:"products" validate:"required"`
}
