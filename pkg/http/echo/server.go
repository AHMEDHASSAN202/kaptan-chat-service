package echo

import (
	"context"
	"example.com/fxdemo/pkg/config"
	echomiddleware "example.com/fxdemo/pkg/http/echo/middleware"
	echoserver "example.com/fxdemo/pkg/http/echo/server"
	"example.com/fxdemo/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
)

func RunServers(lc fx.Lifecycle, log logger.ILogger, e *echo.Echo, ctx context.Context, cfg *config.Config) error {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoserver.RunHttpServer(ctx, e, log, &cfg.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()
			e.Use(echomiddleware.AppendLangMiddleware)
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "working fine")
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
