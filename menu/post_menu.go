package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	emoji "github.com/jayco/go-emoji-flag"
	"libretaxi/context"
	"libretaxi/objects"
	"libretaxi/rabbit"
	"libretaxi/validation"
	"log"
	"strings"
)

type PostMenuHandler struct {
	user *objects.User
	context *context.Context
}

func (handler *PostMenuHandler) postToAdminChannel(text string) {
	msg := tgbotapi.NewMessage(handler.context.Config.Admin_Channel_Chat_Id, text)
	if len(handler.user.Username) == 0 {
		msg.ParseMode = "MarkdownV2"
	}
	banKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("☝️ Shadow ban",fmt.Sprintf("{\"Action\":\"SHADOW_BAN\",\"Id\":%d}", handler.user.UserId)),
		),
	)
	msg.ReplyMarkup = banKeyboard

	handler.context.Send(msg)
}

func (handler *PostMenuHandler) informUsersAround(lon float64, lat float64, text string, postId int64) {
	textWithContacts := ""

	if len(handler.user.Username) == 0 {
		userTextContact := fmt.Sprintf("[%s %s](tg://user?id=%d)", handler.user.FirstName, handler.user.LastName, handler.user.UserId)
		textWithContacts = fmt.Sprintf("%s\n\nvia %s", text, userTextContact)
	} else {
		textWithContacts = fmt.Sprintf("%s\n\nvia @%s", text, handler.user.Username)
	}

	// Post to the admin channel first, do not bother in case of shadow ban
	if !handler.user.ShadowBanned {
		handler.postToAdminChannel(textWithContacts)
	}

	// In case of shadow ban, post to current user only and return
	if handler.user.ShadowBanned {
		msg := tgbotapi.NewMessage(handler.user.UserId, textWithContacts)

		if len(handler.user.Username) == 0 {
			msg.ParseMode = "MarkdownV2"
		}

		handler.context.RabbitPublish.PublishTgMessage(rabbit.MessageBag{
			Message: msg,
			PostId: postId,
			Priority: 3,
		})
		return
	}

	// Post to users around
	userIds := handler.context.Repo.UserIdsAround(lon, lat)

	for i, _ := range userIds {
		userId := userIds[i]
		msg := tgbotapi.NewMessage(userId, textWithContacts)

		if len(handler.user.Username) == 0 {
			// IMPORTANT!
			// Currently markdown not escaped. It means when a user sends, for example "Looking for passenger(s)",
			// Telegram will NOT accept this request, it will fail with
			// Bad Request: can't parse entities: Character '(' is reserved and must be escaped with the preceding '\'
			// TODO: always escape markdown (user input)
			msg.ParseMode = "MarkdownV2"
		}

		reportKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("☝️️Report ⚠️",fmt.Sprintf("{\"Action\":\"REPORT_POST\",\"Id\":%d}", postId)),
			),
		)
		msg.ReplyMarkup = reportKeyboard

		// Mass-send with lower priority (3 instead of 0)
		handler.context.RabbitPublish.PublishTgMessage(rabbit.MessageBag{
			Message: msg,
			PostId: postId,
			Priority: 3,
		})
	}
}

func (handler *PostMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Post menu")

	handler.user = user
	handler.context = context

	if context.Repo.UserPostedRecently(user.UserId) {

		msg := tgbotapi.NewMessage(user.UserId, "🕙 Wait for 5 minutes")
		context.Send(msg)

		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)

	} else if len(message.Text) == 0 {

		flag := ""

		if len(user.LanguageCode) > 0 {
			flag = " " + emoji.GetFlag(user.LanguageCode)
		}

		msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf("Copy & paste text starting with 🚗 (driver) or 👋 (passenger) in the following format (you can use your own language%s), or /cancel, examples:", flag))
		context.Send(msg)

		// IMPORTANT! Do not use Markdown'ish symbols here, like (, ), [, ]... because when user copies and pastes the
		// message below, when username isn't set up, "informUsersAround" method above will turn message into
		// markdown. If these symbols are present, it will mix up the entire message (probably won't be accepted by Telegram)

		msg = tgbotapi.NewMessage(user.UserId, `🚗 Ride offer
From: foobar square
To: airport
Date: today
Time: now
Payment: cash, venmo`)
		context.Send(msg)

		msg = tgbotapi.NewMessage(user.UserId, `👋🏻 Ride wanted
From: foobar st, 42
To: downtown
Date: today
Time: now
Pax: 1`)
		context.Send(msg)
	} else {

		textValidation := validation.NewTextValidation()
		error := textValidation.Validate(message.Text)

		if len(error) > 0 {
			msg := tgbotapi.NewMessage(user.UserId, error + " or /cancel")

			context.Send(msg)
			return
		}

		cleanText := strings.TrimSpace(message.Text)

		post := &objects.Post{
			UserId: user.UserId,
			Text: cleanText,
			Lon: user.Lon,
			Lat: user.Lat,
			ReportCnt: 0,
		}

		context.Repo.SavePost(post);

		handler.informUsersAround(post.Lon, post.Lat, cleanText, post.PostId)

		msg := tgbotapi.NewMessage(user.UserId, "✅ Sent to users around you (25km)")
		context.Send(msg)

		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)
	}
}

func NewPostMenu() *PostMenuHandler {
	handler := &PostMenuHandler{}
	return handler
}
