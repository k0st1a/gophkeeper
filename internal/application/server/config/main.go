// Package config for read config from env and flags.
package config

import (
	"flag"
	"fmt"
	"os"
)

// Config - структура с конфигурационными параметрами сервера.
type Config struct {
	// Address - адрес GRPC эндпоинта сервера в формате host:port (по умолчанию `localhost:8080`).
	// Задается через флаг `-address=<ЗНАЧЕНИЕ>` или переменную окружения `ADDRESS=<ЗНАЧЕНИЕ>`.
	Address string
	// HTTPAddress - адрес HTTP эндпоинта сервера в формате host:port
	// Задается через флаг `-http-address=<ЗНАЧЕНИЕ>` или переменную окружения `HTTP_ADDRESS=<ЗНАЧЕНИЕ>`.
	HTTPAddress string
	// LogLevel - уровень логирования. Возможные значения: debug, info, warn, error (по умолчанию info).
	// Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения `LOG_LEVEL=<ЗНАЧЕНИЕ>`.
	LogLevel string
	// DatabaseURI - адрес подключения к базе данных.
	// Задается через флаг `-dsn=<ЗНАЧЕНИЕ>` или переменную окружения `DATABASE_DSN=<ЗНАЧЕНИЕ>`.
	DatabaseURI string
	// SecretKey - ключ с помощью которого шифруются/проверяются пароли пользователя при регистрации и логине.
	// Задается через флаг `-secret-key=<ЗНАЧЕНИЕ>` или переменную окружения `SECRET_KEY=<ЗНАЧЕНИЕ>`.
	SecretKey string
}

var (
	defaultAddress     = "localhost:8080"
	defaultHTTPAddress = ""
	defaultLogLevel    = "info"
)

// New - создать конфигурацию сервера из аргументов командой строки и переменных окружения.
func New() (*Config, error) {
	cfg := Config{
		Address:     defaultAddress,
		HTTPAddress: defaultHTTPAddress,
		LogLevel:    defaultLogLevel,
	}

	err := cfg.applyFromEnvAndArgs()
	if err != nil {
		return nil, fmt.Errorf("apply config from env and args:%w", err)
	}

	return &cfg, nil
}

func (c *Config) applyFromEnvAndArgs() error {
	// From ENV
	a, ok := os.LookupEnv("ADDRESS")
	if ok {
		c.Address = a
	}

	ha, ok := os.LookupEnv("HTTP_ADDRESS")
	if ok {
		c.HTTPAddress = ha
	}

	ll, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		c.LogLevel = ll
	}

	du, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		c.DatabaseURI = du
	}

	sk, ok := os.LookupEnv("SECRET_KEY")
	if ok {
		c.SecretKey = sk
	}

	flag.StringVar(&c.Address, "address", c.Address, "GRPC endpoint сервера в формате host:port.")
	flag.StringVar(&c.HTTPAddress, "http-address", c.HTTPAddress, "HTTP endpoint сервера в формате host:port.")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel,
		"Уровень логирования. Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения "+
			"`LOG_LEVEL=<ЗНАЧЕНИЕ>.\nВозможные значения: debug, info, warn, error.")
	flag.StringVar(&c.DatabaseURI, "dsn", c.DatabaseURI,
		"Адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -dsn")
	flag.StringVar(&c.SecretKey, "secret-key", c.SecretKey,
		"Ключ, с помощью которого шифруются/проверяются пароли пользователя при регистрации и логине."+
			"Задается через флаг `-secret-key=<ЗНАЧЕНИЕ>` или переменную окружения `SECRET_KEY=<ЗНАЧЕНИЕ>`")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args:%v", flag.Args())
	}

	return nil
}
