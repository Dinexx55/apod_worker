package service

import (
	"encoding/json"
	"fmt"
	"nasa-apod-app/internal/config"
	"nasa-apod-app/internal/models"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type APODWorker struct {
	ApodService    *ApodImagesService
	APIKey         string
	ApodURL        string
	RunImmediately bool
	RunTime        time.Time
	Logger         *zap.Logger
}

func NewAPODWorker(apodService *ApodImagesService, apiKey string, workerConfig config.WorkerConfig, logger *zap.Logger) *APODWorker {
	return &APODWorker{
		ApodService:    apodService,
		APIKey:         apiKey,
		ApodURL:        workerConfig.ApiURL,
		RunImmediately: workerConfig.RunFetchingOnStart,
		RunTime:        workerConfig.RunTime,
		Logger:         logger,
	}
}

func (w *APODWorker) Start() {
	if w.RunImmediately {
		go w.fetchAPOD()
	}

	go func() {
		for {
			nextRun := w.calculateNextRunTime(time.Now())
			w.Logger.Info("Next APOD fetch scheduled at", zap.String("time", nextRun.Format(time.RFC3339)))

			time.Sleep(time.Until(nextRun))

			go w.fetchAPOD()
		}
	}()
}

func (w *APODWorker) calculateNextRunTime(currentTime time.Time) time.Time {
	nextRun := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), w.RunTime.Hour(), w.RunTime.Minute(), 0, 0, currentTime.Location())

	if currentTime.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	return nextRun
}

func (w *APODWorker) fetchAPOD() {
	w.Logger.Info("Fetching APOD data from NASA API with URL: " + w.ApodURL)

	url := fmt.Sprintf("%s?api_key=%s", w.ApodURL, w.APIKey)
	w.Logger.Info("URL: " + url)

	resp, err := http.Get(url)
	if err != nil {
		w.Logger.Error("Failed to request APOD API", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.Logger.Error("APOD API returned non-200 status", zap.Int("status_code", resp.StatusCode))
		return
	}

	var apodData models.APODResponse
	if err := json.NewDecoder(resp.Body).Decode(&apodData); err != nil {
		w.Logger.Error("Failed to decode APOD response", zap.Error(err))
		return
	}

	err = w.ApodService.SaveAPODData(apodData)
	if err != nil {
		w.Logger.Error("Failed to save APOD data", zap.Error(err))
		return
	}

	w.Logger.Info("APOD data successfully fetched and saved.")
}
