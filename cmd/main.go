package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/amqp"
	"github.com/SamuelLeutner/wpp-golang-broaker/internal/whatsapp"
	"github.com/SamuelLeutner/wpp-golang-broaker/router"
	"github.com/gofiber/fiber/v2"
	amqplib "github.com/rabbitmq/amqp091-go"
)

func main() {
	app := fiber.New()
	router.SetupRoutes(app)

	user, password := os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD")
	host, port := os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT")

	if user == "" || password == "" || host == "" || port == "" {
		log.Fatal("Variáveis de ambiente RABBITMQ_USER, RABBITMQ_PASSWORD, RABBITMQ_HOST e RABBITMQ_PORT devem ser definidas")
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	conn, err := amqplib.Dial(url)
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
