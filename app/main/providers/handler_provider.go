package providers

import (
	notificationhandler "github.com/brunobotter/notification-system/api/notification_handler"
	"github.com/brunobotter/notification-system/api/websocket_handler"
	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/brunobotter/notification-system/main/container"
)

type HandlerServiceProvider struct{}

func NewHandlerServiceProvider() *HandlerServiceProvider {
	return &HandlerServiceProvider{}
}

func (p *HandlerServiceProvider) Register(c container.Container) {
	c.Singleton(func(hub web_socket.Hub, logger logger.Logger) *websocket_handler.WebSocketHandler {
		return websocket_handler.NewWebSocketHandler(hub, logger)
	})
	c.Singleton(func(logger logger.Logger) *notificationhandler.NotificationHandler {
		return notificationhandler.NewNotificationHandler(logger)
	})
}
