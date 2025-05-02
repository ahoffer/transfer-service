package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/transfer-service/internal/config"
	"github.com/transfer-service/internal/handler"
	"github.com/transfer-service/internal/repository"
	"github.com/transfer-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting transfer service on port %d...", cfg.Port)

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBConfig.Host,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Name,
		cfg.DBConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	jobRepo := repository.NewJobRepository(db)
	requestRepo := repository.NewJobRequestRepository(db)

	// Migrate both repositories
	if err := jobRepo.Migrate(); err != nil {
		log.Fatalf("Failed to migrate job table: %v", err)
	}
	if err := requestRepo.Migrate(); err != nil {
		log.Fatalf("Failed to migrate job request table: %v", err)
	}

	jobService := service.NewJobService(jobRepo, requestRepo)
	jobHandler := handler.NewJobHandler(jobService, "http://localhost:8080/jobs")

	// Initialize Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Service is healthy")
	})

	// Job endpoints
	jobs := r.Group("/jobs")
	{
		jobs.POST("", jobHandler.CreateJob)
		jobs.GET("/:id", jobHandler.GetJob)
		jobs.POST("/:id/cancel", jobHandler.CancelJob)
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
