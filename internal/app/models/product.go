package models

type Product struct {
	Id           int     `json:"id,omitempty"`
	Image       string  `json:"image,omitempty"`
	Name         string  `json:"name" validate:"required"`
	Price        float64 `json:"price"`
	Rating       float32 `json:"rating"`
	Category     string  `json:"category"`
	CountInStock int     `json:"count_in_stock"`
	Description  string  `json:"description"`
	Sales        bool    `json:"sales"`
	SalesPrice   float64 `json:"sales_price"`
	IsFavourite  *bool    `json:"is_favourite,omitempty"`
	DateCreated  string  `json:"date_created,omitempty"`
	BrandName    string  `json:"brand_name,omitempty"`
}

type OnePageProduct struct {
	Id           int     `json:"id,omitempty"`
	Images       []string  `json:"image,omitempty"`
	Name         string  `json:"name" validate:"required"`
	Price        float64 `json:"price"`
	Rating       float32 `json:"rating"`
	Category     string  `json:"category"`
	CountInStock int     `json:"count_in_stock"`
	Description  string  `json:"description"`
	Sales        bool    `json:"sales"`
	SalesPrice   float64 `json:"sales_price"`
	IsFavourite  *bool    `json:"is_favourite,omitempty"`
	DateCreated  string  `json:"date_created,omitempty"`
	BrandName    string  `json:"brand_name"`
}

type ProductsCategory struct {
	Products []Product `json:"products" validate:"required"`
	MinPrice float64   `json:"min_price" validate:"required"`
	MaxPrice float64   `json:"max_price" validate:"required"`
}

type SalesProduct struct {
	ProductId int `json:"product_id" validate:"required"`
	SalesPrice float64 `json:"sales_price" validate:"required"`
}

type ProductPrice struct {
	Id    int
	Price float64
}

type FavouriteProduct struct {
	Id     int `json:"id" validate:"required"`
	UserId int `json:"user_id,-"`
}

type ProductRating struct {
	Rating int
	Count  int
}

type ProductId struct {
	ProductId int `json:"product_id" validate:"required"`
}

type ProductForSuggest struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ProductIdRecommendCount struct {
	Id      int
	Counter int
}
