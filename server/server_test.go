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

func TestCreateProduct(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(productJson))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	product := new(database.Product)
	if err := c.Bind(product); err != nil {
		t.Errorf("Could not bind product from request body: %v", err)
	}
	ctx := context.Background()

	queries := database.New(database.SetupTest())

	insertedProduct, err := queries.CreateProduct(ctx, product.Name, product.Description, product.Price)
	if err != nil {
		t.Errorf("Failed to insert product: %v because: %v", product, err)
	}

	retrievedInsertedProduct, err := queries.GetProduct(ctx, insertedProduct.ID)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, productJson, recorder.Body.String())
		assert.Equal(t, product.ID, retrievedInsertedProduct.ID)
	}
}
