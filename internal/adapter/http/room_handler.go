package http

import (
	"applicationDesignTest/internal/domain"
	roomUseCase "applicationDesignTest/internal/usecase/room"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type RoomHandler struct {
	roomUseCase *roomUseCase.RoomUseCase
}

func NewRoomHandler(uc *roomUseCase.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUseCase: uc}
}

func (h *RoomHandler) GetByHotelID(w http.ResponseWriter, r *http.Request) {
	hotelID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	rooms, err := h.roomUseCase.GetByHotelID(r.Context(), hotelID)
	if err != nil {
		http.Error(w, "Failed to get rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	hotelID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var room domain.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(room.Number) == "" {
		http.Error(w, "Room number is required", http.StatusBadRequest)
		return
	}

	room.HotelID = hotelID
	if err := h.roomUseCase.Create(r.Context(), &room); err != nil {
		switch {
		case strings.Contains(err.Error(), "room number already exists"):
			http.Error(w, "Room number already exists", http.StatusConflict)
		case strings.Contains(err.Error(), "hotel not found"):
			http.Error(w, "Hotel not found", http.StatusNotFound)
		default:
			http.Error(w, "Failed to create room", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}
