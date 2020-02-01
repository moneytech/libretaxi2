package context

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/repository"
	"libretaxi/rabbit"
)

type Context struct {
	Bot  *tgbotapi.BotAPI
	Repo *repository.Repository
	RabbitPublish *rabbit.RabbitClient // for publishing only
}

// drop-in replacement for telegram Send method, posts with highest priority
func (context * Context) Send(message tgbotapi.Chattable) {
	context.RabbitPublish.PublishTgMessage(rabbit.MessageBag{
		Message: message,
		Priority: 0,
	})
}
