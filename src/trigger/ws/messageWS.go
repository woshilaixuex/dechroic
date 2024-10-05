package ws

import (
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/gorilla/websocket"
	"github.com/xinggaoya/qwen-sdk/qwen"
)

var upgrader = websocket.Upgrader{}

type WebSocketHandler struct {
	qwenClient        *qwen.Chat
	messageRepository *infra_repository.MessageRepository
}

func NewWebSocketHandler(apiKey string, messageRepository *infra_repository.MessageRepository) *WebSocketHandler {
	return &WebSocketHandler{
		qwenClient:        qwen.NewWithDefaultChat(apiKey),
		messageRepository: messageRepository,
	}
}
