package providers

import (
	"github.com/brunobotter/notification-system/api/websocket_handler"
	"github.com/brunobotter/notification-system/infra/websocket"
	"github.com/brunobotter/notification-system/main/container"
)

type WebSocketHandlerServiceProvider struct{}

func NewWebSocketHandlerServiceProvider() *WebSocketHandlerServiceProvider {
	return &WebSocketHandlerServiceProvider{}
}

func (p *WebSocketHandlerServiceProvider) Register(c container.Container) {
	c.Singleton(func(hub websocket.Hub) *websocket_handler.WebSocketHandler {
		return websocket_handler.NewWebSocketHandler(hub)
	})
}
