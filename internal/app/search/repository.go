package search

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	GetSuggests(str string) (models.Suggest, error)
	GetSearchResults(searchArray []string) ([][]models.Product, error)
}
