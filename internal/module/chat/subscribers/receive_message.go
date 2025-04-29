package subscribers

import (
	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
	"kaptan/internal/module/chat/domain"
)

var Topics = struct{ ReceiveMessage string }{ReceiveMessage: "receive_message"}

type Subscriber struct {
	chatUseCase domain.ChatUseCase
	eventBus    EventBus.Bus
}

func Message(eventBus EventBus.Bus, chatUseCase domain.ChatUseCase) {
	o := Subscriber{chatUseCase: chatUseCase, eventBus: eventBus}
	err := eventBus.SubscribeAsync(Topics.ReceiveMessage, o.receiveMessage, false)
	if err != nil {
		logrus.Print("MessageSubscriber | ", err)
		return
	}
}

func (o Subscriber) receiveMessage(message interface{}) {

}
