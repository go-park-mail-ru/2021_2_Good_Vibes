package orders

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	PutOrder(order models.Order) (int, error)
	SelectPrices(products []models.OrderProducts) ([]models.ProductPrice, error)
	GetAllOrders(user int) ([]models.Order, error)
	CheckPromoCode(promoCode string) (*models.PromoCode, error)
	GetProductCategory(productId int) (int, error)
	GetOrderById(orderId int) (models.Order, error)
}
