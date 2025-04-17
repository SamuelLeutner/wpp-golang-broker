package main

import (
	"log"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/amqp"
	"github.com/SamuelLeutner/wpp-golang-broaker/internal/whatsapp"
	"github.com/SamuelLeutner/wpp-golang-broaker/router"
	"github.com/gofiber/fiber/v2"
	amqplib "github.com/rabbitmq/amqp091-go"
)

func main() {
	app := fiber.New()
	router.SetupRoutes(app)

	conn, err := amqplib.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatal("Erro ao conectar ao RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Erro ao abrir canal RabbitMQ:", err)
	}
	defer ch.Close()

	amqp.SetMessageSender(&whatsapp.Service{})

	amqp.Channel = ch
	amqp.StartResponseConsumer()

	err = whatsapp.StartClient()
	if err != nil {
		log.Fatal("Erro ao iniciar WhatsApp:", err)
	}

	log.Fatal(app.Listen(":3000"))
	select {}
}
