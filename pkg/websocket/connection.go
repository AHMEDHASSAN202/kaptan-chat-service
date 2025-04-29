package websocket

import (
	"github.com/labstack/echo/v4"
	"kaptan/pkg/logger"
)

func NewConnectionManger(e *echo.Echo, log logger.ILogger) *ChannelManager {
	manager := NewChannelManager()
	go manager.Run()

	e.GET("/ws", func(c echo.Context) error {
		return handleWebSocket(c, manager)
	})

	return manager
}
