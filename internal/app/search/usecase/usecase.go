package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search"
)

type UseCase struct {
	repositorySearch search.Repository
}

func NewSearchUseCase(repositorySearch search.Repository) *UseCase {
	return &UseCase{
		repositorySearch: repositorySearch,
	}
}


func (uc *UseCase) GetSuggests(str string) (models.Suggest, error) {
	suggests, err := uc.repositorySearch.GetSuggests(str)
	if err != nil {
		return models.Suggest{}, err
	}

	return suggests, nil
}

func (uc *UseCase) GetSearchResults(str string) ([]models.Product, error) {
	products, err := uc.repositorySearch.GetSearchResults(str)
	if err != nil {
		return nil, err
	}

	return products, nil
}
