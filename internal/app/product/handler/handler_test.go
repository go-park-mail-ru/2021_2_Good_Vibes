package handler

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/storage/impl"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllProductsSuccessUnit(t *testing.T) {
	var mockStorage = impl.NewStorageProductsMemory()

	mockStorage.AddProduct(product.Product{Id: 1, Image: "images/cat1.jpeg", Name: "cat1", Price: 1000, Rating: 100})
	mockStorage.AddProduct(product.Product{Id: 2, Image: "images/cat2.jpeg", Name: "cat2", Price: 1000, Rating: 100})

	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)
	product2 := product.NewProduct(2, "images/cat2.jpeg", "cat2", 1000, 100)

	var products []product.Product
	products = append(products, product1)
	products = append(products, product2)

	wantedProductResp, _ := json.Marshal(products)

	tests := []struct {
		name string
		wantedJson string
		statusCode int
	} {
		{
			"products",
			string(wantedProductResp) + "\n",
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := echo.New()

			rec, ctx, h := constructRequest("/homepage", router, mockStorage)

			if assert.NoError(t, h.GetAllProducts(ctx)) {
				assert.Equal(t, tt.statusCode, rec.Code)
				assert.Equal(t, tt.wantedJson, rec.Body.String())
			}
		})
	}
}

func constructRequest(target string, router *echo.Echo, mockStorage *impl.StorageProductsMemory) (*httptest.ResponseRecorder, echo.Context, *ProductHandler) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	h := &ProductHandler{mockStorage}
	return rec, ctx, h
}
