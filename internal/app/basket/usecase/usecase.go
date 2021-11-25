package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/basket"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type UseCase struct {
	basketServiceClient proto.BasketServiceClient
}

func NewBasketUseCase(conn *grpc.ClientConn) *UseCase {
	c := proto.NewBasketServiceClient(conn)
	return &UseCase{
		basketServiceClient: c,
	}
}

func (uc *UseCase) PutInBasket(basket models.BasketProduct) error {
	_, err := uc.basketServiceClient.PutInBasket(context.Background(), models.ModelBasketProductToGrpc(basket))
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) GetBasket(userId int) ([]models.BasketProduct, error) {
	products, err := uc.basketServiceClient.GetBasket(context.Background(),
		models.ModelBasketStorageToGrpc(models.BasketStorage{UserId: userId}))
	if err != nil {
		return nil, err
	}

	var productsModel []models.BasketProduct
	for _, element := range products.GetProducts() {
		productsModel = append(productsModel, models.GrpcBasketProductToModel(element))
	}

	return productsModel, nil
}

func (uc *UseCase) DropBasket(userId int) error {
	_, err := uc.basketServiceClient.DropBasket(context.Background(),
		models.ModelBasketStorageToGrpc(models.BasketStorage{UserId: userId}))
	if err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) DeleteProduct(product models.BasketProduct) error {
	_, err := uc.basketServiceClient.DeleteProduct(context.Background(),
		models.ModelBasketProductToGrpc(product))
	if err != nil {
		return err
	}
	return nil
}
