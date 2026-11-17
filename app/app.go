package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "applicationDesignTest/internal/adapter/http"
	"applicationDesignTest/internal/adapter/postgres"
	"applicationDesignTest/internal/usecase"
	bookingUseCase "applicationDesignTest/internal/usecase/booking"
	roomUseCase "applicationDesignTest/internal/usecase/room"
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	// Listen for interrupt and termination signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Println("Server started on :8080")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}

	case <-ctx.Done():
		log.Println("Shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		if err := <-serverErr; err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}
