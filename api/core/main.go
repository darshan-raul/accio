package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/infra", func(c *fiber.Ctx) error {
		// Communicate with Validator and Implementor services
		// You can use HTTP requests or gRPC for communication
		return c.SendString("Infra Service: Request Routed")
	})

	app.Listen(":3001")
}
