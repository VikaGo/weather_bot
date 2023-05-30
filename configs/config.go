package configs

type Config struct {
	TelegramToken       string `env:"TOKEN"`
	WeatherApi          string `env:"WeatherApi"`
	CageGeocodingAPIKey string `env:"CageGeocodingAPIKey"`
}
