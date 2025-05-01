package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Service is healthy")
	})

	// Handle job creation
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		jobHandler.CreateJob(w, r)
	})

	// Handle job operations (get, cancel)
	http.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Handle cancel endpoint
		if len(pathParts) == 4 && pathParts[3] == "cancel" {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			jobHandler.CancelJob(w, r)
			return
		}

		// Handle get job
		if r.Method == http.MethodGet {
			jobHandler.GetJob(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
