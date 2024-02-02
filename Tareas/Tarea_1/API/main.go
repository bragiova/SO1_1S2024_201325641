package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/data", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"name":   "Brayan Rivas",
			"carnet": 201325641,
		})
	})

	app.Listen(":3000")
	fmt.Println("Server on port 3000")
}
