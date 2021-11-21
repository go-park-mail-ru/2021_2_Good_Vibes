package models

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NestingCategory struct {
	Nesting     int    `json:"nesting"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryNode struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Nesting     int            `json:"-"`
	Children    []CategoryNode `json:"children,omitempty"`
}

type CreateCategory struct {
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Description    string `json:"description"`
}

type CategoryForSuggest struct {
	Name    string `json:"name"`
	Description string `json:"description"`
}
