package sender

import (
	"libretaxi/context"
	"libretaxi/rabbit"
	"log"
)

type Sender struct {
	context *context.Context
}

func (s *Sender) Handler(messageBag *rabbit.MessageBag) {
	log.Printf("Sending message %+v\n", messageBag.Message)
	s.context.Bot.Send(messageBag.Message)
}

func (s *Sender) Start() {
	log.Println("Sender started")
	s.context.RabbitConsume.RegisterHandler(s.Handler)
}

func NewSender(context *context.Context) *Sender {
	sender := &Sender{
		context: context,
	}
	return sender
}
