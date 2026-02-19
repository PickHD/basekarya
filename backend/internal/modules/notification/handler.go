package notification

import (
	"hris-backend/internal/infrastructure"
	"hris-backend/pkg/logger"
	"hris-backend/pkg/response"
	"hris-backend/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	wsHub   *infrastructure.Hub
	service Service
}

func NewHandler(wsHub *infrastructure.Hub, service Service) *Handler {
	return &Handler{wsHub, service}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		logger.Errorw("[WS] Failed to get user context", err)
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	client := &infrastructure.Client{Hub: h.wsHub, UserID: userContext.UserID, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()

	return nil
}

func (h *Handler) GetAll(ctx echo.Context) error {
	userContext, err := utils.GetUserContext(ctx)
	if err != nil {
		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	data, err := h.service.GetList(ctx.Request().Context(), userContext.UserID)
	if err != nil {
		logger.Errorw("get notifications failed: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Get Notification List Success", data, nil, nil)
}

func (h *Handler) MarkAsRead(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.service.MarkAsRead(ctx.Request().Context(), uint(id))
	if err != nil {
		logger.Errorw("failed to mark as read notification: ", err)

		return response.NewResponses[any](ctx, http.StatusInternalServerError, "Failed to mark as read", err.Error(), err, nil)
	}

	return response.NewResponses[any](ctx, http.StatusOK, "Mark as read notification successfully", nil, nil, nil)
}
