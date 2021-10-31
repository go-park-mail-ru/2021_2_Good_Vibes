package models

type Product struct {
	Id           int     `json:"id,omitempty"`
	Image        string  `json:"image"`
	Name         string  `json:"name" validate:"required"`
	Price        float64 `json:"price"`
	Rating       float32 `json:"rating"`
	Category     string  `json:"category"`
	CountInStock int     `json:"count_in_stock"`
	Description  string  `json:"description"`
}

type ProductPrice struct {
	Id    int
	Price float64
}
