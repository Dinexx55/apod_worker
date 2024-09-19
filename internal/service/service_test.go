package service

import (
	"context"
	"errors"
	"nasa-apod-app/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetImageByDate(t *testing.T) {
	logger := zap.NewNop()
	repo := NewInMemoryApodImagesRepo()
	apodService := NewApodImagesService(logger, repo)

	date := "2023-09-18"
	image := domain.ApodImageMetaData{Date: date, Title: "Test"}

	repo.Save(image)

	t.Run("successfully retrieves APOD image by date", func(t *testing.T) {
		img, err := apodService.GetImageByDate(context.Background(), date)
		assert.NoError(t, err)
		assert.Equal(t, &image, img)
	})

	t.Run("invalid date format", func(t *testing.T) {
		_, err := apodService.GetImageByDate(context.Background(), "invalid-date")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDate, err)
	})

	t.Run("image not found", func(t *testing.T) {
		_, err := apodService.GetImageByDate(context.Background(), "2023-09-19")
		assert.Error(t, err)
		assert.Equal(t, ErrImageNotFound, err)
	})
}

func TestGetAllImages(t *testing.T) {
	logger := zap.NewNop()
	repo := NewInMemoryApodImagesRepo()
	apodService := NewApodImagesService(logger, repo)

	image := domain.ApodImageMetaData{Date: "2023-09-18", Title: "Test"}
	repo.Save(image)

	t.Run("successfully retrieves all APOD images", func(t *testing.T) {
		images, err := apodService.GetAllImages(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(images))
		assert.Equal(t, "2023-09-18", images[0].Date)
	})

	t.Run("no images found", func(t *testing.T) {
		repo.images = make(map[string]domain.ApodImageMetaData) // Clear the repo
		_, err := apodService.GetAllImages(context.Background())
		assert.Error(t, err)
		assert.Equal(t, ErrImagesNotFound, err)
	})
}

type InMemoryApodImagesRepo struct {
	images map[string]domain.ApodImageMetaData
}

func NewInMemoryApodImagesRepo() *InMemoryApodImagesRepo {
	return &InMemoryApodImagesRepo{
		images: make(map[string]domain.ApodImageMetaData),
	}
}

func (repo *InMemoryApodImagesRepo) GetAllImages(ctx context.Context) ([]domain.ApodImageMetaData, error) {
	if len(repo.images) == 0 {
		return nil, errors.New("images not found")
	}
	var images []domain.ApodImageMetaData
	for _, img := range repo.images {
		images = append(images, img)
	}
	return images, nil
}

func (repo *InMemoryApodImagesRepo) GetImageByDate(ctx context.Context, date string) (*domain.ApodImageMetaData, error) {
	img, exists := repo.images[date]
	if !exists {
		return nil, errors.New("image not found")
	}
	return &img, nil
}

func (repo *InMemoryApodImagesRepo) ExistsByDate(date string) (bool, error) {
	_, exists := repo.images[date]
	return exists, nil
}

func (repo *InMemoryApodImagesRepo) Save(metadata domain.ApodImageMetaData) error {
	repo.images[metadata.Date] = metadata
	return nil
}
