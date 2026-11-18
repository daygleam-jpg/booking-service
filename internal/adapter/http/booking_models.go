package http

type CreateBookingRequest struct {
	UserID   int    `json:"user_id"`
	HotelID  int    `json:"hotel_id"`
	RoomID   int    `json:"room_id"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}

type CreateBookingResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
