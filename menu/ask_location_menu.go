package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/objects"
	"libretaxi/context"
	"log"
)

type AskLocationMenuHandler struct {
	user *objects.User
	context *context.Context
	message *tgbotapi.Message
}

func (handler *AskLocationMenuHandler) saveLocation() {
	handler.user.Lat = handler.message.Location.Latitude
	handler.user.Lon = handler.message.Location.Longitude
	handler.context.Repo.SaveUser(handler.user)
}

func (handler *AskLocationMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Ask location menu")

	handler.user = user
	handler.context = context
	handler.message = message

	if message.Location != nil {
		log.Printf("Saving location: %+v\n", message.Location)
		handler.saveLocation()

		// update state
		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)
		return
	} else {
		var buttons = []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonLocation("Next"),
		}

		msg := tgbotapi.NewMessage(user.UserId, "Click \"Next\" (from mobile phone) to share your location.")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)
		context.Bot.Send(msg)
	}
}

func NewAskLocationMenu() *AskLocationMenuHandler {
	handler := &AskLocationMenuHandler{}
	return handler
}
