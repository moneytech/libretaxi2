package main

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq" // important
	"libretaxi/callback"
	"libretaxi/config"
	"libretaxi/context"
	"libretaxi/menu"
	"libretaxi/rabbit"
	"libretaxi/repository"
	"libretaxi/sender"
	"log"
)

func initContext() *context.Context {
	log.Printf("Using '%s' telegram token\n", config.C().Telegram_Token)
	log.Printf("Using '%s' database connection string", config.C().Db_Conn_Str)
	log.Printf("Using '%s' RabbitMQ connection string", config.C().Rabbit_Url)

	context := &context.Context{}

	bot, err := tgbotapi.NewBotAPI(config.C().Telegram_Token)
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	db, err := sql.Open("postgres", config.C().Db_Conn_Str)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Successfully connected to the database")
	}

	context.Bot = bot
	context.Repo = repository.NewRepository(db)
	return context
}

// Message producer (app logic)
func main1() {
	context := initContext()
	context.RabbitPublish = rabbit.NewRabbitClient(config.C().Rabbit_Url, "messages")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := context.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			// Ignore the case where message comes from a chat, not from user. We do not support chats.
			if update.Message.From == nil {
				continue
			}

			userId := update.Message.Chat.ID

			// Send welcome link: on first interaction (here) or by mass-send (not implemented yet), but only once per user
			if context.Repo.ShowCallout(userId, "welcome_2_0_message") {
				context.Repo.DismissCallout(userId, "welcome_2_0_message")

				context.Send(tgbotapi.NewMessage(userId, "https://telegra.ph/LibreTaxi-20---you-will-love-it-02-02"))
			}

			log.Printf("[%d - %s] %s", userId, update.Message.From.UserName, update.Message.Text)
			menu.HandleMessage(context, userId, update.Message)

		} else if update.CallbackQuery != nil {

			cb := update.CallbackQuery
			context.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(cb.ID, "👌 Reported"))

			emptyKeyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
			removeButton := tgbotapi.NewEditMessageReplyMarkup(cb.Message.Chat.ID, cb.Message.MessageID, emptyKeyboard)

			_, err := context.Bot.Send(removeButton)
			if err != nil {
				log.Println(err)
			}

			callback.NewTgCallbackHandler().Handle(context, cb.Data)
		}

		//msg := tgbotapi.NewMessage(userId, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//context.Bot.Send(msg)
	}
}

// Message consumer (send to Telegram respecting rate limits)
func main2() {
	context := initContext()
	context.RabbitConsume = rabbit.NewRabbitClient(config.C().Rabbit_Url, "messages")

	s := sender.NewSender(context)
	s.Start()
}

func main() {
	config.Init("libretaxi")

	go main1()
	go main2()

	forever := make(chan bool)
	<- forever
}
