package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/transfer-service/internal/config"
	"github.com/transfer-service/internal/handler"
	"github.com/transfer-service/internal/service"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting transfer service on port %d...", cfg.Port)
	jobService := service.NewJobService()
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
