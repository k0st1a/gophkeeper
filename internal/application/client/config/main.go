// Package config for read config from env and flags.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config - структура с конфигурационными параметрами клиента.
type Config struct {
	// Address - адрес GRPC эндпоинта сервера в формате host:port (по умолчанию `localhost:8080`).
	// Задается через флаг `-address=<ЗНАЧЕНИЕ>` или переменную окружения `ADDRESS=<ЗНАЧЕНИЕ>`.
	Address string
	// LogLevel - уровень логирования. Возможные значения: debug, info, warn, error (по умолчанию info).
	// Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения `LOG_LEVEL=<ЗНАЧЕНИЕ>`.
	LogLevel string
	// LogFile - имя файл, куда будут писаться логи. По умолчанию не задан, в это случае логи пишутся в stdout.
	// Задается через флаг `-log-file=<ЗНАЧЕНИЕ>` или переменную окружения `LOG_FILE=<ЗНАЧЕНИЕ>`.
	LogFile string
	// SecretKey - ключ с помощью которого шифруются/проверяются пароли пользователя при регистрации и логине.
	// Задается через флаг `-secret-key=<ЗНАЧЕНИЕ>` или переменную окружения `SECRET_KEY=<ЗНАЧЕНИЕ>`.
	SecretKey string
	// RequestTimeout - таймаут обращения к серверу, в секундах.
	// Задается через флаг `-request-timeout=<ЗНАЧЕНИЕ>` или переменную окружения `REQUEST_TIMEOUT=<ЗНАЧЕНИЕ>`.
	RequestTimeout int
	// SyncInterval - интервал синхронизации предметов между локальным и удаленным хранилищем, в секундах.
	// Задается через флаг `-sync-interval=<ЗНАЧЕНИЕ>` или переменную окружения `SYNC_INTERVAL=<ЗНАЧЕНИЕ>`.
	SyncInterval int
}

var (
	defaultAddress        = "localhost:8080"
	defaultLogLevel       = "info"
	defaultRequestTimeout = 3
	defaultSyncInterval   = 10
)

// New - создать конфигурацию клиента из аргументов командой строки и переменных окружения.
func New() (*Config, error) {
	cfg := Config{
		Address:        defaultAddress,
		LogLevel:       defaultLogLevel,
		RequestTimeout: defaultRequestTimeout,
		SyncInterval:   defaultSyncInterval,
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

	lf, ok := os.LookupEnv("LOG_FILE")
	if ok {
		c.LogFile = lf
	}

	sk, ok := os.LookupEnv("SECRET_KEY")
	if ok {
		c.SecretKey = sk
	}

	rt, ok := os.LookupEnv("REQUEST_TIMEOUT")
	if ok {
		rtInt, err := strconv.Atoi(rt)
		if err != nil {
			return fmt.Errorf("REQUEST_TIMEOUT parse error:%w", err)
		}
		c.RequestTimeout = rtInt
	}

	si, ok := os.LookupEnv("REQUEST_TIMEOUT")
	if ok {
		siInt, err := strconv.Atoi(si)
		if err != nil {
			return fmt.Errorf("SYNC_INTERVAL parse error:%w", err)
		}
		c.SyncInterval = siInt
	}

	flag.StringVar(&c.Address, "address", c.Address, "GRPC endpoint сервера в формате host:port.")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel,
		"Уровень логирования. Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения "+
			"`LOG_LEVEL=<ЗНАЧЕНИЕ>.\nВозможные значения: debug, info, warn, error.")
	flag.StringVar(&c.LogFile, "log-file", c.LogFile,
		"Файл логирования. Задается через флаг `-log-file=<ЗНАЧЕНИЕ>` или переменную окружения "+
			"`LOG_FILE=<ЗНАЧЕНИЕ>.\nПо умолчанию не задан, в этом случае логи пишутся в stdout.")
	flag.StringVar(&c.SecretKey, "secret-key", c.SecretKey,
		"Ключ, с помощью которого шифруются/проверяются пароли пользователя при регистрации и логине."+
			"Задается через флаг `-secret-key=<ЗНАЧЕНИЕ>` или переменную окружения `SECRET_KEY=<ЗНАЧЕНИЕ>`")
	flag.IntVar(&c.RequestTimeout, "request-timeout", c.RequestTimeout,
		"Таймаут обращения к серверу, в секундах.\n"+
			"Задается через флаг `-request-timeout=<ЗНАЧЕНИЕ>` или переменную окружения `REQUEST_TIMEOUT=<ЗНАЧЕНИЕ>`")
	flag.IntVar(&c.SyncInterval, "sync-interval", c.SyncInterval,
		"Интервал синхронизации элементов между локальным хранилищем и удаленным, в секундах.\n"+
			"Задается через флаг `-sync-interval=<ЗНАЧЕНИЕ>` или переменную окружения `SYNC_INTERVAL=<ЗНАЧЕНИЕ>`")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args:%v", flag.Args())
	}

	return nil
}
