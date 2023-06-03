package pkg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const commandStart = "start"
const commandWeather = "weather"
const commandHelp = "help"

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandWeather:
		return b.handleWeatherCommand(message)
	case commandHelp:
		return b.handleHelpCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Weather forecast for any city with maps.\n☀⛅☔ Just enter a city name as 'text' or send as a 'location' to see the weather forecast 😃 👍")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleWeatherCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Just enter a city name as 'text' or send as a 'location' to see the weather forecast☀️ ☔ ")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "I can help you \n/start \n/weather \n/help")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "The command is invalid")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleLocationUpdate(message *tgbotapi.Message) {
	latitude := message.Location.Latitude
	longitude := message.Location.Longitude

	city, err := getCityName(latitude, longitude)
	if err != nil {
		log.Println("Failed to get city name:", err)
		reply := "Sorry, an error occurred while fetching the city name. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		b.bot.Send(msg)
		return
	}

	forecast, err := getWeatherForecastByLocation(latitude, longitude)
	if err != nil {
		log.Println("Failed to fetch weather forecast:", err)
		reply := "Sorry, an error occurred while fetching the weather forecast. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		b.bot.Send(msg)
		return
	}

	sendWeatherForecast(b.bot, message.Chat.ID, city, forecast)
}

func (b *Bot) handleTextUpdate(message *tgbotapi.Message) {
	city := message.Text

	forecast, err := getWeatherForecastByCity(city)
	if err != nil {
		log.Println("Failed to fetch weather forecast:", err)
		reply := "Sorry, an error occurred while fetching the weather forecast. Please try again later."
		msg := tgbotapi.NewMessage(message.Chat.ID, reply)
		b.bot.Send(msg)
		return
	}

	sendWeatherForecast(b.bot, message.Chat.ID, forecast.City, forecast)
}

func sendWeatherForecast(bot *tgbotapi.BotAPI, chatID int64, city string, forecast *weatherForecast) {
	reply := fmt.Sprintf("Weather forecast for 📍%s:\n\n", city)
	for _, weather := range forecast.WeatherList {
		reply += fmt.Sprintf("Date: %s\nTemperature: %.2f °C\nHumidity: %d%%\nDescription: %s\n\n",
			weather.Date, weather.Temperature, weather.Humidity, weather.Description)
	}

	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}
