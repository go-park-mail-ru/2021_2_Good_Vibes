package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type UseCase struct {
	repository category.Repository
}

func NewCategoryUseCase(repositoryCategory category.Repository) *UseCase {
	return &UseCase{
		repository: repositoryCategory,
	}
}

func (uc *UseCase) GetAllCategories() (models.CategoryNode, error) {
	nestingCategories, err := uc.repository.SelectAllCategories()
	if err != nil {
		return models.CategoryNode{}, err
	}

	node := parseCategories(nestingCategories)

	return node, nil
}

func parseCategories(nestingCategories []models.NestingCategory) models.CategoryNode {
	rootNode := models.CategoryNode {
		Name:    nestingCategories[0].Name,
		Nesting: nestingCategories[0].Nesting,
	}
	var nodeStack []models.CategoryNode
	nodeStack = append(nodeStack, rootNode)

	currentNode := rootNode

	for i := 1; i < len(nestingCategories); i++ {
		if nestingCategories[i].Nesting > nestingCategories[i - 1].Nesting {
			currentNode = models.CategoryNode {
				Name: nestingCategories[i].Name,
				Nesting: nestingCategories[i].Nesting,
				Children: nil,
			}
			nodeStack = append(nodeStack, currentNode)
		}  else if nestingCategories[i].Nesting == nestingCategories[i - 1].Nesting {
			currentNode = models.CategoryNode {
				Name: nestingCategories[i].Name,
				Nesting: nestingCategories[i].Nesting,
				Children: nil,
			}
			nodeStack[len(nodeStack) - 2].Children = append(nodeStack[len(nodeStack) - 2].Children, nodeStack[len(nodeStack) - 1])
			nodeStack = nodeStack[:(len(nodeStack) - 1)]
			nodeStack = append(nodeStack, currentNode)
		} else if nestingCategories[i].Nesting < nestingCategories[i - 1].Nesting {
			diff := nestingCategories[i - 1].Nesting - nestingCategories[i].Nesting

			currentNode = models.CategoryNode {
				Name: nestingCategories[i].Name,
				Nesting: nestingCategories[i].Nesting,
				Children: nil,
			}
			for i := 0; i < diff + 1; i++ {
				nodeStack[len(nodeStack) - i - 2].Children = append(nodeStack[len(nodeStack) - i - 2].Children, nodeStack[len(nodeStack) - i - 1])
			}
			for i := 0; i < diff + 1; i++ {
				nodeStack = nodeStack[:(len(nodeStack) - 1)]
			}
			nodeStack = append(nodeStack, currentNode)
		}
	}
	return nodeStack[0]
}
