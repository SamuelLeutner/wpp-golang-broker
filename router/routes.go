package router

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
	v1.Get("/qrcode", HandlerQRCode)
	v1.Get("/disconnect", HandlerDisconnectPhone)
}
