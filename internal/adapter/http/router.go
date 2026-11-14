package http

import "net/http"

func NewRouter(
	hotelHandler *HotelHandler,
	roomHandler *RoomHandler,
	bookingHandler *BookingHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hotels/{id}", hotelHandler.GetByID)
	mux.HandleFunc("POST /hotels", hotelHandler.Create)
	mux.HandleFunc("GET /hotels", hotelHandler.GetAll)
	mux.HandleFunc("PUT /hotels/{id}", hotelHandler.Update)
	mux.HandleFunc("DELETE /hotels/{id}", hotelHandler.Delete)

	mux.HandleFunc("GET /hotels/{id}/rooms", roomHandler.GetByHotelID)
	mux.HandleFunc("POST /hotels/{id}/rooms", roomHandler.Create)
	mux.HandleFunc("POST /bookings", bookingHandler.Create)
	mux.HandleFunc("DELETE /bookings/{id}", bookingHandler.Cancel)

	return mux
}
