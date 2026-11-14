package domain

import (
	"testing"
	"time"
)

func TestBooking_Validate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		booking   Booking
		wantError bool
	}{
		{
			name: "valid booking",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: false,
		},

		{
			name: "invalid user ID",
			booking: Booking{
				UserID:   0,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{
			name: "invalid hotel ID",
			booking: Booking{
				UserID:   1,
				HotelID:  0,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{
			name: "invalid room ID",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   0,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{
			name: "invalid status",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   "pending",
			},
			wantError: true,
		},
		{

			name: "invalid chekIn",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, -5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{

			name: "invalid checkIn",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, -5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{

			name: "checkOut in day of checkIn",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 5),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
		{

			name: "checkOut before checkIn",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 4),
				Status:   StatusConfirmed,
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.booking.Validate()

			if tt.wantError && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestBooking_Cancel(t *testing.T) {
	now := time.Now()
	cancelled := []struct {
		name      string
		booking   Booking
		wantError bool
	}{
		{
			name: "valid booking",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusConfirmed,
			},
			wantError: false,
		},

		{
			name: "booking already cancelled",
			booking: Booking{
				UserID:   1,
				HotelID:  1,
				RoomID:   1,
				CheckIn:  now.AddDate(0, 0, 5),
				CheckOut: now.AddDate(0, 0, 10),
				Status:   StatusCancelled,
			},
			wantError: true,
		},
	}
	for _, tt := range cancelled {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.booking.Cancel()

			if tt.wantError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantError && tt.booking.Status != StatusCancelled {
				t.Errorf("expected status cancelled, got %v", tt.booking.Status)
			}
		})
	}
}
