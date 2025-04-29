package main

import (
	"go.uber.org/fx"
	"kaptan/internal/module/chat"
	"kaptan/internal/module/user"
	"kaptan/pkg/aws"
	"kaptan/pkg/config"
	"kaptan/pkg/database"
	"kaptan/pkg/event_bus"
	"kaptan/pkg/gate"
	"kaptan/pkg/http"
	"kaptan/pkg/http/echo"
	"kaptan/pkg/http/echo/server"
	"kaptan/pkg/http_client"
	"kaptan/pkg/localization"
	"kaptan/pkg/logger"
	"kaptan/pkg/middlewares"
	"kaptan/pkg/validators"
	"kaptan/pkg/websocket"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				config.Init,
				logger.InitLogger,
				http.NewContext,
				echoserver.NewEchoServer,
				httpclient.NewHttpClient,
				validators.Init,
				aws.Init,
				event_bus.NewEventBus,
			),
			database.Module,
			websocket.Module,
			middlewares.Module,
			gate.Module,
			user.Module,
			chat.Module,
			fx.Invoke(echo.RunServers, localization.InitLocalization),
		),
	).Run()
}
