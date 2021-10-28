package models

type OrderProducts struct {
	OrderId   int `json:"order_id,omitempty"`
	ProductId int `json:"product_id"`
	Number    int `json:"number"`
}

type Address struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Street  string `json:"street"`
	House   string `json:"house"`
	Flat    string `json:"flat"`
	Index   string `json:"index"`
}

type Order struct {
	OrderId  int             `json:"order_id,omitempty"`
	UserId   int             `json:"user_id,omitempty"`
	Date     string          `json:"date"`
	Address  Address         `json:"address"`
	Cost     float32         `json:"cost"`
	Status   string          `json:"status,omitempty"`
	Products []OrderProducts `json:"products"`
}
