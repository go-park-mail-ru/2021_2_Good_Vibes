package models

type OrderProducts struct {
	OrderId   int `json:"order_id,omitempty" validate:"required"`
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
	OrderId  int             `json:"order_id,omitempty" validate:"required"`
	UserId   int             `json:"user_id,omitempty" validate:"required"`
	Date     string          `json:"date" validate:"required"`
	Address  Address         `json:"address" validate:"required"`
	Cost     float32         `json:"cost" validate:"required"`
	Status   string          `json:"status,omitempty" validate:"required"`
	Products []OrderProducts `json:"products" validate:"required"`
}
