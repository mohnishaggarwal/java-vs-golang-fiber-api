package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohnishaggarwal/products/database"
	"github.com/mohnishaggarwal/products/handlers"
	"github.com/mohnishaggarwal/products/repository"
	"log"
)

func main() {
	log.Println("Starting application")
	database.InitDynamo()
	app := fiber.New()

	productRepo := repository.NewProductRepository()
	productHandler := handlers.NewProductHandler(productRepo)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("health check received")
	})
	app.Get("/api/products/:id", productHandler.GetProduct)
	app.Post("/api/products", productHandler.CreateProduct)
	app.Put("/api/products/:id", productHandler.UpdateProduct)

	app.Listen(":8080")
}
