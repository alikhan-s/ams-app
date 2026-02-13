package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"airport-system/internal/airportops"
	"airport-system/internal/auth"
	"airport-system/internal/booking"
	"airport-system/internal/flight"
	"airport-system/internal/passenger"
	"airport-system/platform/database"
	"airport-system/platform/logger"
	"airport-system/platform/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		// Ignore error if file not found
	}

	// 2. Initialize Logger
	log := logger.New()

	// 3. Connect to Database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=airport_db port=5432 sslmode=disable"
	}
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("Connected to database")

	// 4. Initialize TxManager
	txManager := database.NewTxManager(db)
	_ = txManager

	// 5. Setup Gin
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.CORS())

	// 6. Config
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-it"
	}

	// --- MODULE: AUTH ---
	// Wiring dependencies: Repo -> Service -> Handler
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, log, jwtSecret)
	authHandler := auth.NewHandler(authService)

	// API Group
	v1 := router.Group("/api/v1")
	{
		// Register Auth Routes
		auth.RegisterRoutes(v1, authHandler)

		// Register Flight Routes
		flightRepo := flight.NewRepository(db)
		flightService := flight.NewService(flightRepo, log)
		flightHandler := flight.NewHandler(flightService)
		flight.RegisterRoutes(v1, flightHandler)

		// Register Airport Ops Routes
		opsRepo := airportops.NewRepository(db)
		opsService := airportops.NewService(opsRepo, log)
		opsHandler := airportops.NewHandler(opsService)
		authMiddleware := auth.AuthMiddleware(jwtSecret)
		airportops.RegisterRoutes(v1, opsHandler, authMiddleware)

		// Register Passenger Module (Internal dependency, no public routes for now)
		passRepo := passenger.NewRepository(db)
		passService := passenger.NewService(passRepo, log)

		// Register Booking Routes
		bookingRepo := booking.NewRepository(db)
		bookingService := booking.NewService(bookingRepo, flightRepo, txManager, opsService, passService, log)
		bookingHandler := booking.NewHandler(bookingService)
		booking.RegisterRoutes(v1, bookingHandler, authMiddleware)
	}

	// 7. Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exiting")
}
