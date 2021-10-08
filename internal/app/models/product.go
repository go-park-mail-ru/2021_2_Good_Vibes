package models

type Product struct {
	Id       int     `json:"id"`
	Image    string  `json:"image"`
	Name     string  `json:"name"`
	Price    int     `json:"price"`
	Rating   float32 `json:"rating"`
	Category string  `json:"category"`
}
