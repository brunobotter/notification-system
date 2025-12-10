package providers

import (
	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/infra/websocket"
	"github.com/brunobotter/notification-system/main/container"
)

type WebSocketServiceProvider struct{}

func NewWebSocketServiceProvider() *WebSocketServiceProvider {
	return &WebSocketServiceProvider{}
}

func (p *WebSocketServiceProvider) Register(c container.Container) {
	c.Singleton(func(logger logger.Logger) websocket.Hub {
		return websocket.NewHub(logger)
	})
}

func (p *WebSocketServiceProvider) Boot(c container.Container) {
	var hub websocket.Hub
	c.Resolve(&hub)
	go hub.Run()
}
