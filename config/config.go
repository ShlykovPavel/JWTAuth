package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

// Config содержит параметры аутентификации.
// Все переменные обязательны, если не указано env-default.
//
// AUTH_LOGIN_URL    - URL для получения токена
// AUTH_REFRESH_URL  - URL для обновления токена
// AUTH_USERNAME     - Логин сервисного аккаунта
// AUTH_PASSWORD     - Пароль (не логируйте!)
type Config struct {
	Env        string `env:"ENV" env-default:"production"`
	LoginURL   string `env:"AUTH_LOGIN_URL" env-required:"true"`
	RefreshURL string `env:"AUTH_REFRESH_URL" env-required:"true"`
	Username   string `env:"AUTH_USERNAME" env-required:"true"`
	Password   string `env:"AUTH_PASSWORD" env-required:"true"`
	Timeout    int    `env:"AUTH_TIMEOUT_MS" env-default:"5000"`
	RetryCount int    `env:"AUTH_RETRY_COUNT" env-default:"3"`
}

// LoadConfig Загрузка конфига
//
// Сначала загружается конфиг из окружения, если в окружении нет нужных переменных, то происходит попытка загрузки конфига из .env файла
func LoadConfig() *Config {
	var cfg Config
	//Чтение переменных окружения
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Default().Println("Failed to load config from env:", err)
		log.Default().Println("Using default config from .env")
	}

	//	Чтение .env файла
	err = cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalln("Failed to load config from .env file:", err)
	}
	return &cfg
}
