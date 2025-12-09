package providers

import (
	"github.com/brunobotter/notification-system/infra/websocket"
	"github.com/brunobotter/notification-system/main/container"
)

type WebSocketServiceProvider struct{}

func NewWebSocketServiceProvider() *WebSocketServiceProvider {
	return &WebSocketServiceProvider{}
}

func (p *WebSocketServiceProvider) Register(c container.Container) {
	c.Singleton(func() websocket.Hub {
		return websocket.NewHub()
	})
}
