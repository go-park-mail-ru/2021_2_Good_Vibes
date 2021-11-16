package helpers

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

func ParseCategories(nestingCategories []models.NestingCategory) models.CategoryNode {
	if nestingCategories == nil {
		return models.CategoryNode{}
	}

	rootNode := models.CategoryNode{
		Name:    nestingCategories[0].Name,
		Description: nestingCategories[0].Description,
		Nesting: nestingCategories[0].Nesting,
	}
	var nodeSlice []models.CategoryNode
	nodeSlice = append(nodeSlice, rootNode)

	currentNode := rootNode

	for i := 1; i < len(nestingCategories); i++ {
		if nestingCategories[i].Nesting > nestingCategories[i-1].Nesting {
			currentNode = models.CategoryNode{
				Name:     nestingCategories[i].Name,
				Description: nestingCategories[i].Description,
				Nesting:  nestingCategories[i].Nesting,
				Children: nil,
			}
			nodeSlice = append(nodeSlice, currentNode)
		} else if nestingCategories[i].Nesting == nestingCategories[i-1].Nesting {
			currentNode = models.CategoryNode{
				Name:     nestingCategories[i].Name,
				Description: nestingCategories[i].Description,
				Nesting:  nestingCategories[i].Nesting,
				Children: nil,
			}
			nodeSlice[len(nodeSlice)-2].Children = append(nodeSlice[len(nodeSlice)-2].Children, nodeSlice[len(nodeSlice)-1])
			nodeSlice = nodeSlice[:(len(nodeSlice) - 1)]
			nodeSlice = append(nodeSlice, currentNode)
		} else if nestingCategories[i].Nesting < nestingCategories[i-1].Nesting {
			diff := nestingCategories[i-1].Nesting - nestingCategories[i].Nesting

			currentNode = models.CategoryNode{
				Name:     nestingCategories[i].Name,
				Description: nestingCategories[i].Description,
				Nesting:  nestingCategories[i].Nesting,
				Children: nil,
			}
			for i := 0; i < diff+1; i++ {
				nodeSlice[len(nodeSlice)-i-2].Children = append(nodeSlice[len(nodeSlice)-i-2].Children, nodeSlice[len(nodeSlice)-i-1])
			}
			for i := 0; i < diff+1; i++ {
				nodeSlice = nodeSlice[:(len(nodeSlice) - 1)]
			}
			nodeSlice = append(nodeSlice, currentNode)
		}

		if i == len(nestingCategories)-1 {
			for i := 0; i < len(nodeSlice)-1; i++ {
				nodeSlice[len(nodeSlice)-i-2].Children = append(nodeSlice[len(nodeSlice)-i-2].Children, nodeSlice[len(nodeSlice)-i-1])
			}
		}
	}
	return nodeSlice[0]
}
