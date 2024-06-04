package http

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("context is canceled!")
				cancel()
				return
			}
		}
	}()

	return ctx
}
