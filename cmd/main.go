package main

import (
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"nasa-apod-app/internal/app"
	"nasa-apod-app/internal/config"
	"time"
)

func main() {
	envFilePath := flag.String("env", ".env", "path to .env file")
	flag.Parse()

	err := godotenv.Load(*envFilePath)
	if err != nil {
		log.Printf("launching without .env file: %v", err)
	}

	appConfig, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Panic(err)
	}

	logger, err := initLogger()
	if err != nil {
		log.Panicf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	a, err := app.NewApp(appConfig, logger)
	if err != nil {
		log.Panic(err)
	}

	if err := a.Run(); err != nil {
		log.Panic(err)
	}

}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
