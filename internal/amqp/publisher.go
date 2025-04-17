package amqp

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Channel *amqp.Channel

func PublishIncomingMessage(sender string, msg string) error {
	body, err := json.Marshal(map[string]string{
		"from":    sender,
		"message": msg,
	})
	if err != nil {
		return err
	}

	return Channel.Publish(
		"",             // exchange
		"to_wabot_api", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
