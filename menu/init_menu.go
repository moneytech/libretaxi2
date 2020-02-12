package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/context"
	"libretaxi/objects"
	"log"
)

type InitMenuHandler struct {
	context *context.Context
}

func (handler *InitMenuHandler) postToPublicChannel(user *objects.User) {
	text := ""

	if len(user.Username) == 0 {
		userTextContact := fmt.Sprintf("[%s %s](tg://user?id=%d)", user.FirstName, user.LastName, user.UserId)
		text = fmt.Sprintf("%s has joined LibreTaxi", userTextContact)
	} else {
		text = fmt.Sprintf("@%s has joined LibreTaxi", user.Username)
	}

	msg := tgbotapi.NewMessage(handler.context.Config.Public_Channel_Chat_Id, text)
	if len(user.Username) == 0 {
		msg.ParseMode = "MarkdownV2"
	}
	handler.context.Send(msg)
}

func (handler *InitMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Init menu")

	handler.context = context

	// Send welcome message
	msg := tgbotapi.NewMessage(user.UserId, user.Locale().Get("init_menu.welcome"))
	context.Send(msg)

	user.MenuId = objects.Menu_AskLocation
	context.Repo.SaveUser(user)

	handler.postToPublicChannel(user);
}

func NewInitMenu() *InitMenuHandler {
	handler := &InitMenuHandler{}
	return handler
}
