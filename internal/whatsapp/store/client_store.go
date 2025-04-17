package store

import (
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func GetClient() *whatsmeow.Client {
	if client == nil {
		return client
	}

	dbLog := waLog.Stdout("Database", "ERROR", false)
	container, err := sqlstore.New("sqlite3", "file:database.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatal(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatal(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", false)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	return client
}
