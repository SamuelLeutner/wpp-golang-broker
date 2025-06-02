package whatsapp

import (
	"log"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/amqp"
	"go.mau.fi/whatsmeow/types/events"
)

func HandleWhatsAppEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsGroup {
			msg := ""
			if v.Message.GetConversation() != "" {
				msg = v.Message.GetConversation()
			} else if v.Message.GetExtendedTextMessage() != nil {
				msg = v.Message.GetExtendedTextMessage().GetText()
			}

			sender := v.Info.Sender.User

			log.Printf("Nova mensagem de %s: %s", sender, msg)

			err := amqp.PublishIncomingMessage(sender, msg)
			if err != nil {
				log.Println("Erro ao publicar no RabbitMQ:", err)
			}
		}
	}
}
