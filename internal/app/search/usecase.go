package search

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	GetSuggests(str string) (models.Suggest, error)
	GetSearchResults(searchString []string) ([]models.Product, error)
}
