package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr        string
	RedisPassword    string
	KafkaBrokers     []string
	ReaderServerPort int
	WriterServerPort int
	CacheCleanupCron string
	CacheCapacity    int
	KafkaTopic       string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	kafkaBrokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	readerPort, _ := strconv.Atoi(getEnv("READ_SERVER_PORT", "8080"))
	writerPort, _ := strconv.Atoi(getEnv("WRITE_SERVER_PORT", "8081"))
	cacheCleanupCron := getEnv("CACHE_CLEANUP_CRON", "0 0 * * *")
	cacheCapacity, _ := strconv.Atoi(getEnv("CACHE_CAPACITY", "10000"))
	kafkaTopic := getEnv("KAFKA_TOPIC", "cache-updates")

	var kafkaBrokers []string
	if kafkaBrokersStr != "" {
		kafkaBrokers = []string{kafkaBrokersStr}
	}

	return &Config{
		RedisAddr:        redisAddr,
		RedisPassword:    redisPassword,
		KafkaBrokers:     kafkaBrokers,
		ReaderServerPort: readerPort,
		WriterServerPort: writerPort,
		CacheCleanupCron: cacheCleanupCron,
		CacheCapacity:    cacheCapacity,
		KafkaTopic:       kafkaTopic,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
