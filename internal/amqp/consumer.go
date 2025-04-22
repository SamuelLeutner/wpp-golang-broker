package amqp

import (
	"encoding/json"
	"log"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/core"
)

var sender core.MessageSender

func SetMessageSender(ms core.MessageSender) {
	sender = ms
}

func StartResponseConsumer() {
	_, err := Channel.QueueDeclare(
		"to_wabot_api",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	msgs, err := Channel.Consume("to_wabot_api", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			var payload struct {
				To      string `json:"to"`
				Message string `json:"message"`
			}

			if err := json.Unmarshal(d.Body, &payload); err != nil {
				continue
			}

			sender.SendMessage(payload.To, payload.Message)
			log.Printf("Mensagem enviada para %s: %s", payload.To, payload.Message)
		}
	}()
}
