package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"nasa-apod-app/internal/domain"
)

type ApodImagesRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *ApodImagesRepository {
	repo := &ApodImagesRepository{
		db: db,
	}
	return repo
}

func (r *ApodImagesRepository) Save(metadata domain.ApodImageMetaData) error {
	query := `
       INSERT INTO apod_images (title, explanation, date, local_storage_path, copyright)
       VALUES ($1, $2, $3, $4, $5)
   `
	_, err := r.db.Exec(query, metadata.Title, metadata.Explanation, metadata.Date, metadata.LocalStorageImagePath, metadata.Copyright)
	if err != nil {
		return fmt.Errorf("failed to save APOD data: %w", err)
	}
	return nil
}

func (r *ApodImagesRepository) GetImageByDate(ctx context.Context, date string) (*domain.ApodImageMetaData, error) {
	query := `
       SELECT id, title, explanation, date, local_storage_path, copyright
       FROM apod_images
       WHERE date = $1
   `

	var imageMeta domain.ApodImageMetaData
	err := r.db.GetContext(ctx, &imageMeta, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get APOD image by date: %w", err)
	}

	return &imageMeta, nil
}

func (r *ApodImagesRepository) GetAllImages(ctx context.Context) ([]domain.ApodImageMetaData, error) {
	query := `
       SELECT id, title, explanation, date, local_storage_path, copyright
       FROM apod_images
   `

	var images []domain.ApodImageMetaData
	err := r.db.SelectContext(ctx, &images, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all APOD images: %w", err)
	}

	return images, nil
}

func (r *ApodImagesRepository) ExistsByDate(date string) (bool, error) {
	var count int
	query := "SELECT COUNT(1) FROM apod_images WHERE date = $1"

	err := r.db.Get(&count, query, date)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
