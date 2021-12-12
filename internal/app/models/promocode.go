package models

const (
	TypePromoInterest = 1
	TypePromoMoney = 2
)
type PromoCode struct {
	Type int
	Code string
	Value int
	CategoryId int
	ProductId int
	UsesLeft int
}
