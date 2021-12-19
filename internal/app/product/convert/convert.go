package convert

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

func FromProductToOnePageProduct(product models.Product) models.OnePageProduct {
	var outProduct models.OnePageProduct
	outProduct.Id = product.Id
	outProduct.Name = product.Name
	outProduct.Price = product.Price
	outProduct.Rating = product.Rating
	outProduct.Sales = product.Sales
	outProduct.SalesPrice = product.SalesPrice
	outProduct.Description = product.Description
	outProduct.Category = product.Category
	outProduct.CountInStock = product.CountInStock
	outProduct.IsFavourite = product.IsFavourite
	outProduct.BrandName = product.BrandName
	outProduct.DateCreated = product.DateCreated

	return outProduct
}
