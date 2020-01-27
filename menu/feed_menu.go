package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/objects"
	"libretaxi/context"
	"log"
)

type FeedMenuHandler struct {
}

func getKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔍 Search"),
			tgbotapi.NewKeyboardButtonLocation("🌎 Set location"),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func (handler *FeedMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Feed menu")

	if len(message.Text) == 0 && message.Location == nil {

		msg := tgbotapi.NewMessage(user.UserId, "You'll see new posts here. Use 🔍 to search for a 🚗 driver or 🤵 passenger.")
		msg.ReplyMarkup = getKeyboard()
		context.Bot.Send(msg)

	} else if message.Text == "🔍 Search" {

		user.MenuId = objects.Menu_Post
		context.Repo.SaveUser(user)

	} else if message.Location != nil {

		user.Lon = message.Location.Longitude
		user.Lat = message.Location.Latitude
		context.Repo.SaveUser(user)

		msg := tgbotapi.NewMessage(user.UserId, "👌 Location updated")
		msg.ReplyMarkup = getKeyboard()
		context.Bot.Send(msg)

	} else {

		msg := tgbotapi.NewMessage(user.UserId, "😕 Can't understand your choice")
		msg.ReplyMarkup = getKeyboard()
		context.Bot.Send(msg)

	}

	//user.MenuId = objects.Menu_AskLocation
	//context.Repo.SaveUser(user)
}

func NewFeedMenu() *FeedMenuHandler {
	handler := &FeedMenuHandler{}
	return handler
}
