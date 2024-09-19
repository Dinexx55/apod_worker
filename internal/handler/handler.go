package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"nasa-apod-app/internal/domain"
	"nasa-apod-app/internal/service"
	"net/http"
)

type APODImagesService interface {
	GetAllImages(ctx context.Context) ([]domain.ApodImageMetaData, error)
	GetImageByDate(ctx context.Context, date string) (*domain.ApodImageMetaData, error)
}

type APODImagesHandler struct {
	apodService APODImagesService
	logger      *zap.Logger
}

func (h *APODImagesHandler) Init(r *mux.Router) {
	r.HandleFunc("/api/apod", h.GetAllImages).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/apod/{date}", h.GetImageByDate).Methods(http.MethodOptions, http.MethodGet)
}

func NewApodImagesHandler(apodService APODImagesService, logger *zap.Logger) *APODImagesHandler {
	return &APODImagesHandler{
		apodService: apodService,
		logger:      logger,
	}
}

func (h *APODImagesHandler) GetAllImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.apodService.GetAllImages(r.Context())
	if err != nil {
		h.logger.Error("failed to get images", zap.Error(err))

		if errors.Is(err, service.ErrImagesNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "Images not found")
			return
		}

		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func (h *APODImagesHandler) GetImageByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	image, err := h.apodService.GetImageByDate(r.Context(), date)
	if err != nil {
		h.logger.Error("failed to get image", zap.Error(err))

		if errors.Is(err, service.ErrInvalidDate) {
			writeErrorResponse(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD.")
			return
		}

		if errors.Is(err, service.ErrImageNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "Image not found")
			return
		}

		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

func writeErrorResponse(w http.ResponseWriter, code int, message string) {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
