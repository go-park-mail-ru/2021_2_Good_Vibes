package models

type Brand struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name"`
	Image string `json:"image,omitempty"`
}
