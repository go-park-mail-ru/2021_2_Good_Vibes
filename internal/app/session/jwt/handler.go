package jwt

import (
	"context"
)

//go:generate mockgen -source=handler.go -destination=mocks/manager_mock.go
type TokenManager interface {
	GetToken(id int, name string) (string, error)
	ParseTokenFromContext(ctx context.Context) (uint64, error)
}
