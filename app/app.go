package app

import (
	"fmt"
	"log"
	"net/http"

	httpadapter "applicationDesignTest/internal/adapter/http"
	"applicationDesignTest/internal/adapter/postgres"
	"applicationDesignTest/internal/usecase"
	bookingUseCase "applicationDesignTest/internal/usecase/booking"
	roomUseCase "applicationDesignTest/internal/usecase/room"
	"database/sql"
)

type App struct {
	db     *sql.DB
	router http.Handler
}

func New() (*App, error) {
	db, err := postgres.NewDB()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

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

	router := httpadapter.NewRouter(
		hotelHandler,
		roomHandler,
		bookingHandler,
	)

	return &App{
		db:     db,
		router: router,
	}, nil
}

func (a *App) Run() error {
	defer a.db.Close()

	log.Println("Server started on :8080")

	if err := http.ListenAndServe(":8080", a.router); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
