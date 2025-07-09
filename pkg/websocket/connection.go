package websocket

import (
	"github.com/labstack/echo/v4"
	"kaptan/internal/module/chat/domain"
	"kaptan/pkg/logger"
	usermiddleware "kaptan/pkg/middlewares/user"
)

func NewConnectionManger(e *echo.Echo, log logger.ILogger, userMiddleware *usermiddleware.Middlewares, chatUseCase domain.ChatRepository) *ChannelManager {
	manager := NewChannelManager()
	go manager.Run()

	driverChat := e.Group("driver/ws")
	driverChat.Use(userMiddleware.AuthenticationMiddleware("driver"))
	driverChat.GET("", func(c echo.Context) error {
		return handleWebSocket(c, manager, chatUseCase)
	})

	// Important: listen on port 39678
	//e.Logger.Fatal(e.Start(":39678"))

	return manager
}
