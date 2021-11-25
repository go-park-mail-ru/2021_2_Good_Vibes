package models

type Suggest struct {
	Products   []ProductForSuggest  `json:"products"`
	Categories []CategoryForSuggest `json:"categories"`
}
