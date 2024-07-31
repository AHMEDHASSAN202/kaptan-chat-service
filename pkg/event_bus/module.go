package event_bus

import (
	"github.com/asaskevich/EventBus"
	"go.uber.org/fx"
)

// EventBusModule provides an instance of EventBus
var EventBusModule = fx.Module(
	"eventbus",
	fx.Provide(NewEventBus),
)

// NewEventBus creates a new instance of EventBus
func NewEventBus() EventBus.Bus {
	return EventBus.New()
}
