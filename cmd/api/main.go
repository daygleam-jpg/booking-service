package main

import (
	"log"
	"net/http"

	httpadapter "applicationDesignTest/internal/adapter/http"
	"applicationDesignTest/internal/adapter/postgres"
	"applicationDesignTest/internal/usecase"
	bookingUseCase "applicationDesignTest/internal/usecase/booking"
	roomUseCase "applicationDesignTest/internal/usecase/room"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	// Hotels
	hotelRepo := postgres.NewHotelRepository(db)
	hotelUC := usecase.NewHotelUseCase(hotelRepo)
	hotelHandler := httpadapter.NewHotelHandler(hotelUC)

	// Rooms
	roomRepo := postgres.NewRoomRepository(db)
	roomUC := roomUseCase.NewRoomUseCase(roomRepo)
	roomHandler := httpadapter.NewRoomHandler(roomUC)

	// Bookings
	bookingRepo := postgres.NewBookingRepository(db)
	bookingUC := bookingUseCase.NewBookingUseCase(roomRepo, bookingRepo)
	bookingHandler := httpadapter.NewBookingHandler(bookingUC)

	router := httpadapter.NewRouter(hotelHandler, roomHandler, bookingHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
