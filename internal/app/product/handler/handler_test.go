package handler

import (
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

	mockStorage.AddProduct(product.Product{"1.jpg", "cat1", 10000})
	mockStorage.AddProduct(product.Product{"2.jpg", "cat2", 10000})

	tests := []struct {
		name string
		wantedJson string
		statusCode int
	} {
		{
			"products",
			`[{"image":"1.jpg","name":"cat1","price":10000},{"image":"2.jpg","name":"cat2","price":10000}]` + "\n",
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
