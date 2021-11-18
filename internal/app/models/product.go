package models

type Product struct {
	Id           int     `json:"id,omitempty"`
	Image        string  `json:"image,omitempty"`
	Name         string  `json:"name"`
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

type FavouriteProduct struct {
	Id int `json:"id" validate:"required"`
	UserId int `json:"user_id,-"`
}

type ProductRating struct {
	Rating int
	Count int
}

type ProductId struct {
	ProductId int `json:"product_id" validate:"required"`
}
