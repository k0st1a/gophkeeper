// Package config for read config from env and flags.
package config

import (
	"flag"
	"fmt"
	"os"
)

// Config - структура с конфигурационными параметрами сервера.
type Config struct {
	// Address - адрес эндпоинта сервера (по умолчанию `localhost:8080`).
	// Задается через флаг `-a=<ЗНАЧЕНИЕ>` или переменную окружения `ADDRESS=<ЗНАЧЕНИЕ>`
	Address string
	// LogLevel - уровень логирования. Возможные значения: debug, info, warn, error (по умолчанию info).
	// Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения `LOG_LEVEL=<ЗНАЧЕНИЕ>`
	LogLevel string
}

var (
	defaultAddress  = "localhost:8080"
	defaultLogLevel = "info"
)

// New - создать конфигурацию сервера из аргументов командой строки и переменных окружения.
func New() (*Config, error) {
	cfg := Config{
		Address:  defaultAddress,
		LogLevel: defaultLogLevel,
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

	ll, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		c.LogLevel = ll
	}

	flag.StringVar(&c.Address, "address", c.Address, "Endpoint сервера в формате host:port.")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel,
		"Уровень логирования. Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения "+
			"`LOG_LEVEL=<ЗНАЧЕНИЕ>.\nВозможные значения: debug, info, warn, error.")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args:%v", flag.Args())
	}

	return nil
}
