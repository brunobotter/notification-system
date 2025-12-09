package websocket_handler

import (
	"net/http"

	"github.com/brunobotter/notification-system/infra/websocket"
	socket "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// WebSocketHandler lida com as conexões WebSocket.
type WebSocketHandler struct {
	hub websocket.Hub
}

// NewWebSocketHandler cria um novo handler de WebSocket.
func NewWebSocketHandler(hub websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

var upgrader = socket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Ajuste para produção!
}

// Handle faz o upgrade da conexão HTTP para WebSocket e registra o cliente no Hub.
func (h *WebSocketHandler) Handle(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := websocket.NewClient(conn, h.hub)
	h.hub.Register(client)

	go client.ReadPump()
	go client.WritePump()

	return nil // A conexão WebSocket é mantida aberta pelos pumps
}
