package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/service"
)

type JobHandler struct {
	jobService *service.JobService
	baseURL    string
}

func NewJobHandler(jobService *service.JobService, baseURL string) *JobHandler {
	return &JobHandler{
		jobService: jobService,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
	}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.jobService.CreateJob(&req)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Set Location header with full URL
	location := fmt.Sprintf("%s/jobs/%d", h.baseURL, job.ID)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract job ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	jobID, err := strconv.ParseUint(pathParts[2], 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := h.jobService.GetJob(uint(jobID))
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) CancelJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract job ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	jobID, err := strconv.ParseUint(pathParts[2], 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	err = h.jobService.CancelJob(uint(jobID))
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
