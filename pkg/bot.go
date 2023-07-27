package pkg

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}
type Database struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type Subscription struct {
	ChatID   int64  `bson:"_id"`
	Time     string `bson:"time"`
	Location string `bson:"location"`
}

func (b *Bot) Start(database Database) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		} else {
			// Check if the message is a time value for subscription
			if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == b.bot.Self.ID {
				b.subscribeUser(database.Collection, chatID, messageText, "")
				reply := fmt.Sprintf("You have subscribed to daily weather forecast notifications at %s (UTC)!", messageText)
				msg := tgbotapi.NewMessage(chatID, reply)
				b.bot.Send(msg)
			} else {
				b.subscribeUser(database.Collection, chatID, messageText, "")
				reply := fmt.Sprintf("You have subscribed to daily weather forecast notifications for %s!", messageText)
				msg := tgbotapi.NewMessage(chatID, reply)
				b.bot.Send(msg)
			}
		}
		continue
	}
	return nil
}

func (b *Bot) subscribeUser(collection *mongo.Collection, chatID int64, time string, city string) {
	subscription := Subscription{
		ChatID:   chatID,
		Time:     time,
		Location: city,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	_, err := collection.InsertOne(ctx, subscription)
	if err != nil {
		log.Println("Error inserting subscription:", err)
	}
}
