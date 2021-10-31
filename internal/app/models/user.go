package models

import "database/sql"

type UserDataForInput struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserDataStorage struct {
	Id       int            `json:"id"`
	Name     string         `json:"username" validate:"required"`
	Email    string         `json:"email"    validate:"required,email"`
	Password string         `json:"password" validate:"required"`
	Avatar   sql.NullString `json:"avatar"`
}

type UserDataForReg struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,customPassword"`
}

type UserID struct {
	UserId int `json:"user_id"`
}

type UserDataProfile struct {
	Id     uint64         `json:"id,omitempty"`
	Name   string         `json:"username"`
	Email  string         `json:"email" validate:"email"`
	Avatar sql.NullString `json:"avatar"`
}
