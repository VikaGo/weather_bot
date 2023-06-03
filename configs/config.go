package configs

type Config struct {
	TelegramToken       string `env:"TOKEN"`
	WeatherApi          string `env:"WAPI"`
	CageGeocodingAPIKey string `env:"APIKEY"`
}
