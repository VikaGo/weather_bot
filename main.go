package main

import (
	"context"
	"github.com/VikaGo/weather_bot/database"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken       string `env:"TOKEN"`
	WeatherApi          string `env:"WAPI"`
	CageGeocodingAPIKey string `env:"APIKEY"`
}
type Subscription struct {
	ChatID   int64  `bson:"_id"`
	Time     string `bson:"time"`
	Location string `bson:"location"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TOKEN := os.Getenv("TOKEN")
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	db, err := database.NewDatabase(mongoURI, dbName, collectionName)
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	c := cron.New()

	c.AddFunc("@daily", func() {
		sendDailyWeatherForecasts(db.Collection, TOKEN)
	})

	c.Start()

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

}

func sendDailyWeatherForecasts(collection *mongo.Collection, botToken string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	currentTime := now.Format("15:04")

	filter := bson.M{"time": currentTime}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Error finding subscriptions:", err)
		return
	}
	defer cursor.Close(ctx)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println("Error creating Telegram bot:", err)
		return
	}

	for cursor.Next(ctx) {
		var subscription Subscription
		if err := cursor.Decode(&subscription); err != nil {
			log.Println("Error decoding subscription:", err)
			continue
		}

		chatID := subscription.ChatID

		msg := tgbotapi.NewMessage(chatID, "Today's weather forecast: Sunny")
		bot.Send(msg)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error iterating over subscriptions:", err)
	}
}

func SubscribeUser(collection *mongo.Collection, chatID int64, time string, city string) {
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
