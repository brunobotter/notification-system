package providers

import (
	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/brunobotter/notification-system/main/container"
)

type WebSocketServiceProvider struct{}

func NewWebSocketServiceProvider() *WebSocketServiceProvider {
	return &WebSocketServiceProvider{}
}

func (p *WebSocketServiceProvider) Register(c container.Container) {
	c.Singleton(func(logger logger.Logger) web_socket.Hub {
		return web_socket.NewHub(logger)
	})
}

func (p *WebSocketServiceProvider) Boot(c container.Container) {
	var hub web_socket.Hub
	c.Resolve(&hub)
	go hub.Run()
}
