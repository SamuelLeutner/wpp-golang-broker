package core

type MessageSender interface {
	SendMessage(to string, message string)
}
