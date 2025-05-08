package amqp

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Channel *amqp.Channel

// Message represents a standardized message format for the broker
type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	ID      string `json:"id"`
}

var (
	ErrInvalidChannel = errors.New("AMQP channel not initialized")
)

// PublishIncomingMessage safely publishes a message with retry logic
func PublishIncomingMessage(sender, message string) error {
	if Channel == nil {
		return ErrInvalidChannel
	}

	msg := Message{
		From:    sender,
		Content: message,
		ID:      uuid.New().String(), // Generate unique message ID
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = Channel.PublishWithContext(ctx,
		"",             // exchange
		"to_wabot_api", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message %s: %w", msg.ID, err)
	}

	log.Printf("Published message %s from %s", msg.ID, sender)
	return nil
}
