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
	q, err := Channel.QueueDeclare(
		"to_wabot_api",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	msgs, err := Channel.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack (false to handle manually)
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			var payload struct {
				To      string `json:"to"`
				Message string `json:"message"`
			}

			var payload struct {
				To      string `json:"to"`
				Message string `json:"message"`
				ID      string `json:"id"`
			}

			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Failed to unmarshal message (ID: %s): %v", payload.ID, err)
				_ = d.Nack(false, false)
				continue
			}

			if payload.To == "" || payload.Message == "" {
				log.Printf("Invalid message payload: %+v", payload)
				_ = d.Nack(false, false) 
				continue
			}

			log.Printf("Processing message %s for %s", payload.ID, payload.To)
			
			sender.SendMessage(payload.To, payload.Message)
			
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message %s: %v", payload.ID, err)
			}
			log.Printf("Successfully processed message %s for %s", payload.ID, payload.To)
		}
	}()
}
