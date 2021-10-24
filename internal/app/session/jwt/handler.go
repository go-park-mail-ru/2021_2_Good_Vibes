package jwt

import "github.com/labstack/echo/v4"

//go:generate mockgen -source=handler.go -destination=mocks/manager_mock.go
type TokenManager interface {
	GetToken(id int, name string) (string, error)
	ParseTokenFromContext(ctx echo.Context) (uint64, error)
}

