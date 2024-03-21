package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	database "github.com/RobinCombrink/ecommerce_saas/database"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var productJson = `{
	"name": "Wooden Chair",
	"description": "Four legs, solid back",
	"price": 493.29
}`

var productDescription string = "Four legs, solid back"
var productParams database.CreateProductParams = database.CreateProductParams{Name: "Wooden Char", Description: &productDescription, Price: 49.99}

var productIdsJson = `[
	1, 2, 3, 4, 5, 6
]`

func TestCreateProduct(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	product := new(database.CreateProductParams)
	if err := c.Bind(product); err != nil {
		t.Errorf("Could not bind product from request body: %v", err)
	}
	ctx := context.Background()

	queries := database.New(database.SetupTest())

	insertedProduct, err := queries.CreateProduct(ctx, *product)
	if err != nil {
		t.Errorf("Failed to insert product: %v because: %v", product, err)
	}

	retrievedInsertedProduct, err := queries.GetProduct(ctx, insertedProduct.ID)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, insertedProduct.ID, retrievedInsertedProduct.ID)
	}
}

func TestDeleteProducts(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/remove-products", strings.NewReader(productIdsJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	productIds := new([]int64)
	if err := c.Bind(productIds); err != nil {
		t.Errorf("Could not bind productIds from request body: %v", err)
	}
	ctx := context.Background()

	queries := database.New(database.SetupTest())

	for i := 0; i < 6; i++ {
		queries.CreateProduct(ctx, productParams)
	}

	deletedProducts, err := queries.DeleteProducts(ctx, *productIds)
	if err != nil {
		t.Errorf("Failed to delete products: %v because: %v", productIds, err)
	}
	assert.NoError(t, err)
	assert.Equal(t, len(deletedProducts), len(*productIds))
	// TODO
	// retrievedInsertedProduct, err := queries.GetProduct(ctx, deletedProducts.ID)
	// if assert.NoError(t, err) {
	// 	assert.Equal(t, http.StatusOK, recorder.Code)
	// 	assert.Equal(t, deletedProducts, retrievedInsertedProduct.ID)
	// }
}
