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
	Username   string `env:"AUTH_USERNAME" env-required:"true"`
	Password   string `env:"AUTH_PASSWORD" env-required:"true"`
	RetryCount int    `env:"AUTH_RETRY_COUNT" env-default:"3"`
}

// LoadConfig Загрузка конфига
//
// Сначала загружается конфиг из окружения, если в окружении нет нужных переменных, то происходит попытка загрузки конфига из .env файла
func LoadConfig(filePath string) *Config {
	var cfg Config
	//Чтение переменных окружения
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Default().Println("Failed to load config from env:", err)
		log.Default().Println("Using default config from .env")
	} else {
		return &cfg
	}
	//	Чтение .env файла
	err = cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		log.Fatalln("Failed to load config from .env file:", err)
	}
	return &cfg
}

//TODO Добавить чтение конфига из yaml
