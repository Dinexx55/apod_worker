package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"nasa-apod-app/internal/domain"
	"nasa-apod-app/internal/models"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ApodImagesRepo interface {
	GetAllImages(ctx context.Context) ([]domain.ApodImageMetaData, error)
	GetImageByDate(ctx context.Context, date string) (*domain.ApodImageMetaData, error)
	ExistsByDate(date string) (bool, error)
	Save(metadata domain.ApodImageMetaData) error
}

type ApodImagesService struct {
	logger     *zap.Logger
	repository ApodImagesRepo
}

var (
	ErrInvalidDate    = fmt.Errorf("invalid date format provided. use YYYY-MM-DD")
	ErrImageNotFound  = fmt.Errorf("image not found")
	ErrImagesNotFound = fmt.Errorf("images not found")
)

func NewApodImagesService(logger *zap.Logger, repository ApodImagesRepo) *ApodImagesService {
	return &ApodImagesService{
		logger:     logger,
		repository: repository,
	}
}

func (s *ApodImagesService) SaveAPODData(apodData models.APODResponse) error {
	exists, err := s.repository.ExistsByDate(apodData.Date)
	if err != nil {
		s.logger.Error("Failed to check if APOD data exists", zap.Error(err))
		return fmt.Errorf("failed to check if APOD data exists: %w", err)
	}

	if exists {
		s.logger.Info("APOD data already exists", zap.String("date", apodData.Date))
		return fmt.Errorf("APOD for today was already saved")
	}

	s.logger.Info("Saving APOD data", zap.String("date", apodData.Date))

	imagePath, err := s.downloadImage(apodData.URL, apodData.Date)
	if err != nil {
		s.logger.Error("Failed to download image", zap.Error(err))
		return fmt.Errorf("failed to download image: %w", err)
	}

	metadata := domain.ApodImageMetaData{
		Title:                 apodData.Title,
		Explanation:           apodData.Explanation,
		Date:                  apodData.Date,
		Copyright:             apodData.Copyright,
		LocalStorageImagePath: imagePath,
	}

	err = s.repository.Save(metadata)
	if err != nil {
		s.logger.Error("Failed to save APOD data", zap.Error(err))
		return err
	}

	s.logger.Info("APOD data saved successfully", zap.String("date", apodData.Date))
	return nil
}

func (s *ApodImagesService) downloadImage(imageURL, date string) (string, error) {
	s.logger.Info("Downloading image", zap.String("url", imageURL), zap.String("date", date))

	resp, err := http.Get(imageURL)
	if err != nil {
		s.logger.Error("Failed to download image", zap.Error(err))
		return "", fmt.Errorf("failed to download image from %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Non-200 response code received", zap.Int("status_code", resp.StatusCode))
		return "", fmt.Errorf("received non-200 response code while downloading image: %d", resp.StatusCode)
	}

	storageDir := "./storage/apod"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		s.logger.Error("Failed to create storage directory", zap.Error(err))
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	fileName := fmt.Sprintf("%s.jpg", date)
	filePath := filepath.Join(storageDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create image file", zap.Error(err))
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		s.logger.Error("Failed to save image to file", zap.Error(err))
		return "", fmt.Errorf("failed to save image to file: %w", err)
	}

	s.logger.Info("Image saved successfully", zap.String("file_path", filePath))
	return filePath, nil
}

func (s *ApodImagesService) GetImageByDate(ctx context.Context, date string) (*domain.ApodImageMetaData, error) {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		s.logger.Error("Invalid date format", zap.Error(err))
		return nil, ErrInvalidDate
	}

	image, err := s.repository.GetImageByDate(ctx, date)
	if err != nil {
		s.logger.Error("Failed to fetch APOD image by date", zap.Error(err))
		return nil, ErrImageNotFound
	}

	s.logger.Info("APOD image fetched successfully", zap.String("date", date))
	return image, nil
}

func (s *ApodImagesService) GetAllImages(ctx context.Context) ([]domain.ApodImageMetaData, error) {
	images, err := s.repository.GetAllImages(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch all APOD images", zap.Error(err))
		return nil, err
	}

	if len(images) == 0 {
		s.logger.Error("no images found in storage", zap.Error(err))
		return nil, ErrImagesNotFound
	}

	s.logger.Info("All APOD images fetched successfully")
	return images, nil
}
