package notificationhandler

import (
	"net/http"

	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	logger logger.Logger
}

func NewNotificationHandler(logger logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		logger: logger,
	}
}

// DTO de entrada (Request)
type createNotificationRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// POST /notifications
func (h *NotificationHandler) Create(c echo.Context) error {
	var req createNotificationRequest

	// 1️⃣ Bind do JSON
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// 2️⃣ Validação básica
	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user_id is required",
		})
	}

	if req.Message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "message is required",
		})
	}

	// 3️⃣ Chamada do service (ainda não implementado)
	// h.service.Notify(req.UserID, req.Message)

	// 4️⃣ Resposta HTTP
	return c.NoContent(http.StatusCreated)
}
