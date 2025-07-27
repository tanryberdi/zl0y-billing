package main

import (
	"log"

	"zl0y-billing/internal/config"
	"zl0y-billing/internal/database"
	"zl0y-billing/internal/handlers"
	"zl0y-billing/internal/middleware"
	"zl0y-billing/internal/repository"
	"zl0y-billing/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize the database connections
	pgDB, err := database.NewPostgresDB(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgDB.Close()

	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect()

	// initialize repositories
	userRepo := repository.NewUserRepository(pgDB)
	reportRepo := repository.NewReportRepository(mongoDB)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, reportRepo)
	reportService := service.NewReportService(reportRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	reportHandler := handlers.NewReportHandler(reportService)
	mockHandler := handlers.NewMockHandler(reportRepo)

	// Setup routes
	router := gin.Default()

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Mock routes for testing
	mock := router.Group("/api/mock")
	{
		mock.POST("/create-report", mockHandler.CreateReport)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.POST("/user/link-anonymous", userHandler.LinkAnonymous)
		protected.GET("/user/reports", userHandler.GetReports)
		protected.POST("/reports/:report_id/purchase", reportHandler.PurchaseReport)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
