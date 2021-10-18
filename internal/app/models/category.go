package models

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NestingCategory struct {
	Nesting int    `json:"nesting"`
	Name    string `json:"name"`
}

type CategoryNode struct {
	Name     string         `json:"name"`
	Nesting  int            `json:"-"`
	Children []CategoryNode `json:"children,omitempty"`
}
