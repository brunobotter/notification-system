package router

import (
	notificationhandler "github.com/brunobotter/notification-system/api/notification_handler"
	"github.com/brunobotter/notification-system/api/websocket_handler"
	"github.com/brunobotter/notification-system/main/config"
	"github.com/brunobotter/notification-system/main/container"
	"github.com/labstack/echo/v4"
)

func RegisterRouter(e *echo.Echo, cfg *config.Config, c container.Container) {
	// Resolva o handler do container
	var wsHandler *websocket_handler.WebSocketHandler
	var notification *notificationhandler.NotificationHandler
	c.Resolve(&wsHandler)
	c.Resolve(&notification)
	// Registre a rota WebSocket
	e.GET("/ws", wsHandler.WebScoketHandle)
	e.POST("/notification", notification.Create)
}
