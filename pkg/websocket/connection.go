package websocket

import (
	"github.com/labstack/echo/v4"
	"kaptan/pkg/logger"
	usermiddleware "kaptan/pkg/middlewares/user"
)

func NewConnectionManger(e *echo.Echo, log logger.ILogger, userMiddleware *usermiddleware.Middlewares) *ChannelManager {
	manager := NewChannelManager()
	go manager.Run()

	driverChat := e.Group("driver/ws")
	driverChat.Use(userMiddleware.AuthenticationMiddleware("driver"))
	driverChat.GET("", func(c echo.Context) error {
		return handleWebSocket(c, manager)
	})

	// Important: listen on port 39678
	//e.Logger.Fatal(e.Start(":39678"))

	return manager
}
