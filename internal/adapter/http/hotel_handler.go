package http

import (
	"applicationDesignTest/internal/domain"
	"applicationDesignTest/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type HotelHandler struct {
	hotelUseCase *usecase.HotelUseCase
}

func NewHotelHandler(uc *usecase.HotelUseCase) *HotelHandler {
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
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotel); err != nil {
		http.Error(w, "Failed to encode hotel", http.StatusInternalServerError)
		return
	}
}
func (h *HotelHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.hotelUseCase.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to get hotels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotels); err != nil {
		http.Error(w, "Failed to encode hotels", http.StatusInternalServerError)
		return
	}
}

func (h *HotelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var hotel domain.Hotel

	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(hotel.Name) == "" {
		http.Error(w, "Hotel name is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(hotel.Location) == "" {
		http.Error(w, "Hotel name is required", http.StatusBadRequest)
		return
	}
	if err := h.hotelUseCase.Create(r.Context(), &hotel); err != nil {
		http.Error(w, "Failed to create hotel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(hotel); err != nil {
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

	var hotel domain.Hotel

	if err := h.hotelUseCase.Update(r.Context(), &hotel); err != nil {
		if err.Error() == "hotel not found" {
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hotel.ID = id

	if err := h.hotelUseCase.Update(r.Context(), &hotel); err != nil {
		if errors.Is(err, errors.New("hotel not found")) {
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotel); err != nil {
		return
	}
}

func (h *HotelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid hotel ID", http.StatusBadRequest)
		return
	}

	if err := h.hotelUseCase.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete hotel",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
