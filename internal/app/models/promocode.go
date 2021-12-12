package models

type PromoCode struct {
	Type int
	Code string
	Value int
	CategoryId int
	ProductId int
	UsesLeft int
}
