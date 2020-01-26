package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/context"
	"libretaxi/objects"
	"log"
)

type Handler interface {
	Handle(user *objects.User, context *context.Context, message *tgbotapi.Message)
}

func isStateChanged(context *context.Context, previousState objects.MenuId, userId int64) (result bool) {
	user := context.Repo.FindUser(userId)

	if user == nil {
		return true
	}

	return user.MenuId != previousState
}

func HandleMessage(context *context.Context, userId int64, message *tgbotapi.Message) {
	log.Printf("Message: '%s'", message.Text)
	previousState := objects.Menu_Ban

	for isStateChanged(context, previousState, userId) == true {
		user := context.Repo.FindUser(userId)

		// Init user if it's not present
		if user == nil {
			user = &objects.User{
				UserId: userId,
				MenuId: objects.Menu_Init,
			}
		}

		// Save username if it's there and changed
		if message.From != nil && len(message.From.UserName) > 0 && message.From.UserName != user.Username {
			user.Username = message.Chat.UserName
			context.Repo.SaveUser(user)
		}

		//fmt.Printf("%+v\n", message.Location)

		if message.Text == "/start" {
			user.MenuId = objects.Menu_Init
			message.Text = ""
			context.Repo.SaveUser(user)
		}

		if message.Text == "/cancel" {
			user.MenuId = objects.Menu_Feed
			message.Text = ""
			context.Repo.SaveUser(user)
		}

		previousState = user.MenuId
		var handler Handler

		switch user.MenuId {
		case objects.Menu_Init:
			handler = NewInitMenu()
		case objects.Menu_AskLocation:
			handler = NewAskLocationMenu()
		case objects.Menu_Feed:
			handler = NewFeedMenu()
		case objects.Menu_Post:
			handler = NewPostMenu()
		default:
			handler = nil
		}

		if handler != nil {
			handler.Handle(user, context, message)
		} else {
			log.Printf("Handler not implemented for menu with id %d\n", user.MenuId)
		}

		// Important! We need to redefine the message as indicator it has been processed.
		// Otherwise it can go into infinite loop.
		message = &tgbotapi.Message{}
	}
}