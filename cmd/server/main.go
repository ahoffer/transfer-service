package main

import (
	"fmt"
	"log"
	"net/http"

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

	// Initialize repository and service
	jobRepo := repository.NewJobRepository(db)
	if err := jobRepo.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	jobService := service.NewJobService(jobRepo)
	jobHandler := handler.NewJobHandler(jobService, "/jobs")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Service is healthy")
	})

	http.HandleFunc("/jobs", jobHandler.CreateJob)
	http.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			jobHandler.GetJob(w, r)
		case http.MethodDelete:
			jobHandler.CancelJob(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
