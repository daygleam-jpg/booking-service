package http

import (
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type hotelUseCase interface {
	GetByID(ctx context.Context, id int64) (*domain.Hotel, error)
	GetAll(ctx context.Context) ([]domain.Hotel, error)
	Create(ctx context.Context, input usecase.CreateHotelInput) (*domain.Hotel, error)
	Update(ctx context.Context, input usecase.UpdateHotelInput) (*domain.Hotel, error)
	Delete(ctx context.Context, id int64) error
}

type HotelHandler struct {
	hotelUseCase hotelUseCase
}

func NewHotelHandler(uc hotelUseCase) *HotelHandler {
	return &HotelHandler{
		hotelUseCase: uc,
	}
}

func (h *HotelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return

	}

	hotel, err := h.hotelUseCase.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "hotel not found" {
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response := HotelResponse{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Location:    hotel.Location,
		Description: hotel.Description,
		CreatedAt:   hotel.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
func (h *HotelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.hotelUseCase.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get hotels", http.StatusInternalServerError)
		return
	}

	response := make([]HotelResponse, 0, len(hotels))

	for _, hotel := range hotels {
		response = append(response, HotelResponse{
			ID:          hotel.ID,
			Name:        hotel.Name,
			Location:    hotel.Location,
			Description: hotel.Description,
			CreatedAt:   hotel.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *HotelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateHotelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "Hotel name is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Location) == "" {
		http.Error(w, "Hotel location is required", http.StatusBadRequest)
		return
	}

	input := usecase.CreateHotelInput{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	hotel, err := h.hotelUseCase.Create(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to create hotel", http.StatusInternalServerError)
		return
	}

	response := HotelResponse{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Location:    hotel.Location,
		Description: hotel.Description,
		CreatedAt:   hotel.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
func (h *HotelHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var req UpdateHotelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := usecase.UpdateHotelInput{
		ID:          id,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	hotel, err := h.hotelUseCase.Update(r.Context(), input)
	if err != nil {
		if err.Error() == "hotel not found" {
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
		return
	}

	response := HotelResponse{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Location:    hotel.Location,
		Description: hotel.Description,
		CreatedAt:   hotel.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
func (h *HotelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	if err := h.hotelUseCase.Delete(r.Context(), id); err != nil {
		if err.Error() == "hotel not found" {
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
