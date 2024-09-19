package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseConfig DBConfig
	ServerConfig   ServerConfig
	WorkerConfig   WorkerConfig
	NasaApiKey     string
}

type WorkerConfig struct {
	RunTime            time.Time
	RunFetchingOnStart bool
	ApiURL             string
}

type DBConfig struct {
	Host           string
	Port           string
	Username       string
	Password       string
	DBName         string
	ReconnRetry    int
	TimeWaitPerTry time.Duration
}

type ServerConfig struct {
	Host string
	Port string
}

func ParseConfigFromEnv() (*Config, error) {
	dbConfig := DBConfig{
		Host:           getEnvOrDefault("DB_HOST", "localhost"),
		Port:           getEnvOrDefault("DB_PORT", "5432"),
		Username:       getEnvOrDefault("DB_USERNAME", "user"),
		Password:       getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:         getEnvOrDefault("DB_NAME", "database"),
		ReconnRetry:    getEnvAsInt("DB_RECONN_RETRY", 3),
		TimeWaitPerTry: getEnvAsDuration("DB_TIME_WAIT_PER_TRY", 5*time.Second),
	}

	serverConfig := ServerConfig{
		Host: getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
		Port: getEnvOrDefault("SERVER_PORT", "8080"),
	}

	NasaApiKey := os.Getenv("NASA_API_KEY")

	workerConfig := WorkerConfig{
		RunTime:            getEnvAsTime("WORKER_RUN_TIME", "03:00"),
		RunFetchingOnStart: getEnvAsBool("RUN_FETCHING_ON_START", true),
		ApiURL:             getEnvOrDefault("NASA_API_URL", "https://api.nasa.gov/planetary/apod?api_key="+os.Getenv("NASA_API_KEY")),
	}

	return &Config{
		DatabaseConfig: dbConfig,
		ServerConfig:   serverConfig,
		WorkerConfig:   workerConfig,
		NasaApiKey:     NasaApiKey,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsTime(name string, defaultValue string) time.Time {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := time.Parse("15:04", valueStr); err == nil {
			return value
		}
	}
	if defaultTime, err := time.Parse("15:04", defaultValue); err == nil {
		return defaultTime
	}
	return time.Time{}
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
