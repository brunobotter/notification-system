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
	Hub    web_socket.Hub
	logger logger.Logger
}

// NewWebSocketHandler cria um novo handler de WebSocket.
func NewWebSocketHandler(hub web_socket.Hub, logger logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		Hub:    hub,
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

	// Cria o client
	client := web_socket.NewClient(conn, h.Hub, h.logger)

	// Registra no Hub
	h.Hub.Register(client)

	// Inicia o write pump (não bloqueia)
	go client.WritePump()

	// O read pump bloqueia enquanto o WS está vivo
	client.ReadPump()

	// Quando sair do ReadPump: desconectar
	h.Hub.Unregister(client)
	conn.Close()

	return nil
}
