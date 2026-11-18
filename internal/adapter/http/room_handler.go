package http

import (
	"applicationDesignTest/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type RoomHandler struct {
	roomUseCase *usecase.RoomUseCase
}

func NewRoomHandler(uc *usecase.RoomUseCase) *RoomHandler {
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

	response := make([]RoomResponse, 0, len(rooms))

	for _, room := range rooms {
		response = append(response, RoomResponse{
			ID:        room.ID,
			HotelID:   room.HotelID,
			Number:    room.Number,
			Status:    room.Status,
			CreatedAt: room.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	hotelID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Number) == "" {
		http.Error(w, "Room number is required", http.StatusBadRequest)
		return
	}

	input := usecase.CreateRoomInput{
		HotelID: hotelID,
		Number:  req.Number,
		Status:  req.Status,
	}

	room, err := h.roomUseCase.Create(r.Context(), input)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "room number already exists"):
			http.Error(w, "Room number already exists", http.StatusConflict)
		case strings.Contains(err.Error(), "hotel not found"):
			http.Error(w, "Hotel not found", http.StatusNotFound)
		case strings.Contains(err.Error(), "invalid room status"):
			http.Error(w, "Invalid room status", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to create room", http.StatusInternalServerError)
		}
		return
	}

	response := RoomResponse{
		ID:        room.ID,
		HotelID:   room.HotelID,
		Number:    room.Number,
		Status:    room.Status,
		CreatedAt: room.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
