package http

import (
	"applicationDesignTest/internal/usecase"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type BookingHandler struct {
	bookingUC *usecase.BookingUseCase
}

func NewBookingHandler(uc *usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUC: uc}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateBookingInput

	var request struct {
		UserID   int    `json:"user_id"`
		HotelID  int    `json:"hotel_id"`
		RoomID   int    `json:"room_id"`
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input.UserID = request.UserID
	input.HotelID = request.HotelID
	input.RoomID = request.RoomID

	var err error
	input.CheckIn, err = time.Parse("2006-01-02", request.CheckIn)
	if err != nil {
		http.Error(w, "Invalid check_in date", http.StatusBadRequest)
		return
	}

	input.CheckOut, err = time.Parse("2006-01-02", request.CheckOut)
	if err != nil {
		http.Error(w, "Invalid check_out date", http.StatusBadRequest)
		return
	}

	bookingID, err := h.bookingUC.CreateBooking(input)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			http.Error(w, "Room not found", http.StatusNotFound)
		case err.Error() == "room does not belong to the specified hotel",
			err.Error() == "invalid booking identifiers",
			err.Error() == "departure date must be after arrival date",
			err.Error() == "arrival date cannot be in the past":
			http.Error(w, err.Error(), http.StatusBadRequest)
		case err.Error() == "room is under maintenance",
			err.Error() == "room is not available":
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      bookingID,
		"message": "Booking created successfully",
	})
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	if err := h.bookingUC.CancelBooking(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
