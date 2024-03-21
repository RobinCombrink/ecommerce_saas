package server

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "github.com/RobinCombrink/ecommerce_saas/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const serverIp string = "127.0.0.1"
const serverPort string = "5000"
const dateLayout string = "2006-01-02"

// TODO: Make it non global
var queries *database.Queries
var db *sql.DB

func SetupHttpServer() {
	instance := echo.New()
	instance.Pre(middleware.RemoveTrailingSlash())

	instance.Use(middleware.Logger())
	instance.Use(middleware.Recover())

	instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	db = database.Setup()
	queries = database.New(db)

	setupRoutes(instance)
	instance.Logger.Fatal(instance.Start(serverIp + ":" + serverPort))

}

func setupRoutes(instance *echo.Echo) {
	instance.GET("/product/:id", getProduct)
	instance.POST("/product", createProduct)
	instance.POST("/remove-products", deleteProduct)
}

func getProduct(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		echo.Logger.Error(c.Logger(), "Could not parse id %d: %s", id, err)
		//TODO: make error better
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid id %s", idParam))
	}
	ctx := context.Background()
	product, err := queries.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, fmt.Sprintf("No product with id: %d", id))
		}
		echo.Logger.Error(c.Logger(), "Could not get product: %d: %s", id, err)
		return c.JSON(http.StatusInternalServerError, "Internal Server error")
	}
	return c.JSON(http.StatusOK, product)
}

func createProduct(c echo.Context) error {
	product := new(database.CreateProductParams)
	if err := c.Bind(product); err != nil {
		//TODO: Differentiate better
		return c.JSON(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	insertedProduct, err := queries.CreateProduct(ctx, *product)
	if err != nil {
		log.Fatalf("Unable to create product: %v\n", err)
		//TODO
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, insertedProduct)
}

func deleteProduct(c echo.Context) error {
	productIds := make([]int64, 0)
	if err := c.Bind(productIds); err != nil {
		//TODO: Differentiate better
		return c.JSON(http.StatusBadRequest, err)
	}
	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	deletedProducts, err := qtx.DeleteProducts(ctx, productIds)
	if err != nil {
		log.Fatalf("Unable to delete products: %v\n", err)
		//TODO
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Unable to delete products, rolling back: %v\n", err)
		//TODO
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, deletedProducts)
}
