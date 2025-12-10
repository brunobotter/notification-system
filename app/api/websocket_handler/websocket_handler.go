package websocket_handler

import (
	"net/http"

	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/infra/web_socket"
	socket "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocketHandler lida com as conexões WebSocket.
type WebSocketHandler struct {
	hub    web_socket.Hub
	logger logger.Logger
}

// NewWebSocketHandler cria um novo handler de WebSocket.
func NewWebSocketHandler(hub web_socket.Hub, logger logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub:    hub,
		logger: logger,
	}
}

var upgrader = socket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handle faz o upgrade da conexão HTTP para WebSocket e registra o cliente no Hub.
func (h *WebSocketHandler) Handle(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := web_socket.NewClient(conn, h.hub, h.logger)
	h.hub.Register(client)

	go client.ReadPump()
	go client.WritePump()

	return nil // A conexão WebSocket é mantida aberta pelos pumps
}
