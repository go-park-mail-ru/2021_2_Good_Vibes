package models

import (
	"database/sql"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/auth"
)

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
	RealName sql.NullString			`json:"real_name,omitempty"`
	RealSurname sql.NullString		`json:"real_surname,omitempty"`
	Sex	 sql.NullString			`json:"sex,omitempty"`
	BirthDay	sql.NullString	`json:"birth_day,omitempty"`
}

type UserDataForReg struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserID struct {
	UserId int `json:"user_id"`
}

type UserDataProfile struct {
	Id     uint64 `json:"id,omitempty"`
	Name   string `json:"username"`
	Email  string `json:"email" validate:"email"`
	Avatar string `json:"avatar,omitempty"`
	RealName string			`json:"real_name,omitempty"`
	RealSurname string		`json:"real_surname,omitempty"`
	Sex	 string			`json:"sex,omitempty"`
	BirthDay	string	`json:"birthday,omitempty"`
}

type UserDataPassword struct {
	Id       uint64 `json:"id,omitempty"`
	Password string `json:"password"`
}

func GrpcUserDataForInputToModel(grpcData *proto.UserForInput) UserDataForInput {
	return UserDataForInput{
		Name:     grpcData.GetName(),
		Password: grpcData.GetPassword(),
	}
}

func ModelUserDataForInputToGrpc(model UserDataForInput) *proto.UserForInput {
	return &proto.UserForInput{
		Name:     model.Name,
		Password: model.Password,
	}
}

func GrpcUserIdToModel(grpcData *proto.UserId) UserID {
	return UserID{
		UserId: int(grpcData.GetId()),
	}
}

func ModelUserIdToGrpc(modelUserId UserID) *proto.UserId {
	return &proto.UserId{
		Id: int64(modelUserId.UserId),
	}
}

func GrpcUserDataForRegToModel(grpcData *proto.UserForReg) UserDataForReg {
	return UserDataForReg{
		Name:     grpcData.GetName(),
		Password: grpcData.GetPassword(),
		Email:    grpcData.GetEmail(),
	}
}

func ModelUserDataForRegToGrpc(model UserDataForReg) *proto.UserForReg {
	return &proto.UserForReg{
		Name:     model.Name,
		Email:    model.Email,
		Password: model.Password,
	}
}
