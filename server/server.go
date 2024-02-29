package server

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"

	database "github.com/RobinCombrink/ecommerce_saas/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const serverIp string = "127.0.0.1"
const serverPort string = "5000"
const dateLayout string = "2006-01-02"

// TODO: Make it non global
var queries *database.Queries

//go:embed schema.sql
var ddl string

func SetupHttpServer() {
	instance := echo.New()
	instance.Pre(middleware.RemoveTrailingSlash())

	instance.Use(middleware.Logger())
	instance.Use(middleware.Recover())

	instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	// db, err := sql.Open("sqlite3", "database.db")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
	}
	ctx := context.Background()

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		//TODO: Handle better
		log.Printf("Could not create tables: %s", err)
	}

	queries = database.New(db)

	setupRoutes(instance)
	instance.Logger.Fatal(instance.Start(serverIp + ":" + serverPort))

}

func setupRoutes(instance *echo.Echo) {
	instance.GET("/product/:id", getProduct)
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
