package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/objects"
	"libretaxi/context"
	"log"
)

type FeedMenuHandler struct {
}

func (handler *FeedMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Feed menu")

	if len(message.Text) == 0 {

		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("🔍 Search"),
				tgbotapi.NewKeyboardButton("🌎 Set location"),
			),
		)
		keyboard.OneTimeKeyboard = true

		msg := tgbotapi.NewMessage(user.UserId, "You'll see new posts here. Use 🔍 to search for a 🚗 driver or 🤵 passenger.")
		msg.ReplyMarkup = keyboard
		context.Bot.Send(msg)

	} else if message.Text == "🔍 Search" {

		user.MenuId = objects.Menu_Post
		context.Repo.SaveUser(user)

	} else {

		msg := tgbotapi.NewMessage(user.UserId, "😕 Can't understand your choice")
		context.Bot.Send(msg)

	}

	//user.MenuId = objects.Menu_AskLocation
	//context.Repo.SaveUser(user)
}

func NewFeedMenu() *FeedMenuHandler {
	handler := &FeedMenuHandler{}
	return handler
}
