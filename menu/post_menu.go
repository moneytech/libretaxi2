package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/context"
	"libretaxi/objects"
	"libretaxi/validation"
	"log"
	"strings"
)

type PostMenuHandler struct {
}

func (handler *PostMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Post menu")

	if len(message.Text) == 0 {

		msg := tgbotapi.NewMessage(user.UserId, "Send text starting with 🚗 or 👋 in the following format (you can use your own language), or /cancel, examples:")
		context.Bot.Send(msg)

		msg = tgbotapi.NewMessage(user.UserId, `🚗 Driver looking for hitcher
Pick Up: foobar reservoir/nearby
Drop Off: anywhere except town
Date: today
Time: now
Payment: cash, venmo`)
		context.Bot.Send(msg)

		msg = tgbotapi.NewMessage(user.UserId, `👋🏻 Hitcher looking for driver
Pick up: 48a foobar st, Oakland
Drop off: downtown
Date: today
Time: now
Pax: 1`)
		context.Bot.Send(msg)

	} else {

		textValidation := validation.NewTextValidation()
		error := textValidation.Validate(message.Text)

		if len(error) > 0 {
			msg := tgbotapi.NewMessage(user.UserId, error)
			context.Bot.Send(msg)
			return
		}

		post := &objects.Post{
			UserId: user.UserId,
			Text: strings.TrimSpace(message.Text),
			Lat: user.Lat,
			Lon: user.Lon,
		}

		context.Repo.SaveNewPost(post);

		msg := tgbotapi.NewMessage(user.UserId, "✅ Sent to users around you (25km)")
		context.Bot.Send(msg)

		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)
	}
}

func NewPostMenu() *PostMenuHandler {
	handler := &PostMenuHandler{}
	return handler
}
