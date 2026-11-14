package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	bookingUseCase "applicationDesignTest/internal/usecase/booking"
)

type BookingHandler struct {
	bookingUC *bookingUseCase.BookingUseCase
}

func NewBookingHandler(uc *bookingUseCase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUC: uc}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input bookingUseCase.CreateBookingInput

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

	if err := h.bookingUC.CreateBooking(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
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
		http.Error(w, "Failed to cancel booking", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
