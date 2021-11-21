package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/auth"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth"
)

type grpcUserHandler struct {
	userUseCase auth.UseCase
	proto.UnimplementedAuthServiceServer
}

func NewGrpcUserHandler(userUseCase auth.UseCase) *grpcUserHandler {
	return &grpcUserHandler{
		userUseCase: userUseCase,
	}
}

func (handler *grpcUserHandler) Login(ctx context.Context, userInput *proto.UserForInput) (*proto.UserId, error) {
	userID, err := handler.userUseCase.Login(models.GrpcUserDataForInputToModel(userInput))
	if err != nil {
		return models.ModelUserIdToGrpc(models.UserID{UserId: userID}), err
	}
	return models.ModelUserIdToGrpc(models.UserID{UserId: userID}), err
}

func (handler *grpcUserHandler) SignUp(ctx context.Context, userInput *proto.UserForReg) (*proto.UserId, error) {
	userID, err := handler.userUseCase.SignUp(models.GrpcUserDataForRegToModel(userInput))
	if err != nil {
		return models.ModelUserIdToGrpc(models.UserID{UserId: userID}), err
	}
	return models.ModelUserIdToGrpc(models.UserID{UserId: userID}), err
}
