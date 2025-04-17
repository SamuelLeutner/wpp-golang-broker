package whatsapp

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Service struct{}

func (s *Service) SendMessage(to string, message string) {
	SendMessage(to, message)
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func GenerateQRCode() (string, error) {
	dbLog := waLog.Stdout("Database", "ERROR", false)
	container, err := sqlstore.New("sqlite3", "file:database.db?_foreign_keys=on", dbLog)
	if err != nil {
		return "", err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return "", err
	}

	clientLog := waLog.Stdout("Client", "ERROR", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())

		err = client.Connect()
		if err != nil {
			return "", err
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				return evt.Code, nil
			} else if evt.Event == "timeout" || evt.Event == "cancelled" {
				return "", errors.New("Error generating QR Code. Please try again.")
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			return "", err
		}
		return "", errors.New("Already logged in")
	}

	return "", errors.New("unexpected flow")
}

func DisconnectPhone() (string, error) {
	dbLog := waLog.Stdout("Database", "ERROR", false)
	container, err := sqlstore.New("sqlite3", "file:database.db?_foreign_keys=on", dbLog)
	if err != nil {
		return "", err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return "", err
	}

	clientLog := waLog.Stdout("Client", "ERROR", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		return "", errors.New("Client is not logged in")
	}

	err = client.Connect()
	if err != nil {
		return "", fmt.Errorf("error connecting before disconnect: %v", err)
	}

	client.Disconnect()
	return "Disconnected successfully", nil
}
