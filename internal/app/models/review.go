package models

type Review struct {
	UserId    int    `json:"-"`
	UserName  string `json:"username,omitempty"`
	Avatar	string `json:"avatar,omitempty"`
	ProductId int    `json:"product_id,omitempty" validate:"required"`
	Rating    int    `json:"rating" validate:"required"`
	Text      string `json:"text" validate:"required"`
	Date	  string `json:"date,omitempty"`
}
