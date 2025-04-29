package event_bus

import "github.com/asaskevich/EventBus"

func NewEventBus() EventBus.Bus {
	bus := EventBus.New()
	return bus
}
