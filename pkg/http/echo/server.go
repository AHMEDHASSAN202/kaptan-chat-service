package echo

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
	"samm/pkg/config"
	echomiddleware "samm/pkg/http/echo/middleware"
	echoserver "samm/pkg/http/echo/server"
	"samm/pkg/logger"
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
			e.Use(echomiddleware.ServerHeader)

			e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins:  []string{"*"},
				AllowMethods:  []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodPost, http.MethodDelete},
				AllowHeaders:  []string{"*"},
				ExposeHeaders: []string{echomiddleware.HEADER_NAME},
			}))
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
