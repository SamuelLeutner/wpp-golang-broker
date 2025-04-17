package whatsapp

import (
	"context"
	"log"

	"github.com/SamuelLeutner/wpp-golang-broaker/internal/core"
	"github.com/SamuelLeutner/wpp-golang-broaker/internal/whatsapp/store"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func SendMessage(to string, message string) {
	client := store.GetClient()
	jid := types.NewJID(to, "s.whatsapp.net")
	_, err := client.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: &message,
	})
	if err != nil {
		log.Println("Error sending message:", err)
	} else {
		log.Println("Message sent successfully!")
	}
}

var _ core.MessageSender = (*Service)(nil)
