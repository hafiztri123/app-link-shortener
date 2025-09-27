package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const clickQueue = "click_events"

type Config struct {
	DatabaseAddr string
	RedisAddr    string
	IDOffset     uint64
	SecretKey    string
	AppPort      string
	RabbitMQAddr string
	ClicksQueue string
}

func Load() (*Config, error) {
	errors := []string{}

	appUrl := getEnvOrDefault("APP_URL", "localhost")
	idOffset, err := strconv.ParseUint(getEnvOrDefault("ID_OFFSET", "10000000"), 10, 64) 

	if err != nil {
		errors = append(errors, err.Error())
	}


	dbAddr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		getEnvOrDefault("DB_USER", "postgres"), 
		getEnvOrDefault("DB_PASSWORD", "password"), 
		appUrl,
		getEnvOrDefault("DB_PORT", "5432"), 
		getEnvOrDefault("DB_TRANSACTION_NAME", "app_db"),
		getEnvOrDefault("DB_SSL", "disable"),
	)

	redisAddr := fmt.Sprintf("%s:%s", 
		appUrl, 
		getEnvOrDefault("REDIS_PORT", "6379"),
	)

	if len(errors) > 0 {
		slog.Error("Failed to load config", "errors", errors)
		return nil, fmt.Errorf("failed to load config: %s", strings.Join(errors, ", "))
	}


	return &Config{
		DatabaseAddr: dbAddr,
		RedisAddr:    redisAddr,
		IDOffset:     idOffset,
		SecretKey:    getEnvOrDefault("JWT_SECRET", "secret"),
		RabbitMQAddr: getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ClicksQueue: clickQueue,
		AppPort: getEnvOrDefault("APP_PORT", "8080"),
	}, nil

}

func getEnvOrDefault (key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}