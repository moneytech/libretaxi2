package context

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/repository"
)

type Context struct {
	Bot  *tgbotapi.BotAPI
	Repo *repository.Repository
}
