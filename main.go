package main

import (
	"github.com/VikaGo/weather_bot/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	TelegramToken       string `env:"TOKEN"`
	WeatherApi          string `env:"WeatherApi"`
	CageGeocodingAPIKey string `env:"CageGeocodingAPIKey"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TOKEN := os.Getenv("TOKEN")

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	telegramBot := pkg.NewBot(bot)

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
		return
	}

}
