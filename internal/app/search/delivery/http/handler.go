package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search"
	sessionJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
)

type SearchHandler struct {
	useCase        search.UseCase
	sessionManager sessionJwt.TokenManager
}

func NewSearchHandler(useCase search.UseCase, sessionManager sessionJwt.TokenManager) *SearchHandler {
	return &SearchHandler{
		useCase:        useCase,
		sessionManager: sessionManager,
	}
}

const trace = "SearchHandler"

