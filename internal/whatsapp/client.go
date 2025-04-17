package whatsapp

import (
	"github.com/SamuelLeutner/wpp-golang-broaker/internal/whatsapp/store"
)

func StartClient() error {
	client := store.GetClient()
	client.AddEventHandler(HandleWhatsAppEvent)

	err := client.Connect()
	if err != nil {
		return err
	}

	return nil
}
