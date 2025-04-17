package router

import (
	"log"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/whatsapp"
	"github.com/gofiber/fiber/v2"
)

type QrCode struct {
	Code string `json:"code"`
}

func HandlerQRCode(c *fiber.Ctx) error {
	code, err := whatsapp.GenerateQRCode()
	if err != nil {
		log.Println("Error generating QR Code:", err)
		return c.Status(500).SendString("Error generating QR Code")
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "QR Code generated successfully",
		"data":    QrCode{Code: code},
	})
}

func HandlerDisconnectPhone(c *fiber.Ctx) error {
	code, err := whatsapp.DisconnectPhone()
	if err != nil {
		log.Println("Error disconnecting phone:", err)
		return c.Status(500).SendString("Error disconnecting phone")
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Phone disconnected successfully",
		"data":    QrCode{Code: code},
	})
}
