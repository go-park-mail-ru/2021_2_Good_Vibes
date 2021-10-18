package models

type UserDataForInput struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserDataStorage struct {
	Id       int    `json:"id"`
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserDataForReg struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,customPassword"`
}

type UserID struct {
	UserId int `json:"user_id"`
}
