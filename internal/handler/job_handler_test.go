package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/service"
)

func TestJobHandler_CreateJob(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	reqBody := `{
		"name": "Test Job",
		"sourceUrl": "http://source/file.txt",
		"destination": "target-server",
		"destinationUrl": "http://target/receive"
	}`

	req := httptest.NewRequest("POST", "/jobs", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateJob(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if !strings.HasPrefix(location, "http://localhost:8080/jobs/") {
		t.Errorf("Expected Location header to start with 'http://localhost:8080/jobs/', got %q", location)
	}

	// Extract job ID from location
	jobID := strings.TrimPrefix(location, "http://localhost:8080/jobs/")
	if len(jobID) != 16 {
		t.Errorf("Expected job ID length 16, got %d", len(jobID))
	}
}

func TestJobHandler_GetJob(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	// First create a job
	reqBody := `{
		"name": "Test Job",
		"sourceUrl": "http://source/file.txt",
		"destination": "target-server",
		"destinationUrl": "http://target/receive"
	}`

	createReq := httptest.NewRequest("POST", "/jobs", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	handler.CreateJob(createW, createReq)
	location := createW.Result().Header.Get("Location")
	jobID := strings.TrimPrefix(location, "http://localhost:8080/jobs/")

	// Now test getting the job
	getReq := httptest.NewRequest("GET", "/jobs/"+jobID, nil)
	getW := httptest.NewRecorder()
	handler.GetJob(getW, getReq)

	resp := getW.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var job models.Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify job fields
	if job.Name != "Test Job" {
		t.Errorf("Expected job name 'Test Job', got %q", job.Name)
	}
	if job.SourceUrl != "http://source/file.txt" {
		t.Errorf("Expected source URL 'http://source/file.txt', got %q", job.SourceUrl)
	}
	if job.Destination != "target-server" {
		t.Errorf("Expected destination 'target-server', got %q", job.Destination)
	}
	if job.DestinationUrl != "http://target/receive" {
		t.Errorf("Expected destination URL 'http://target/receive', got %q", job.DestinationUrl)
	}
	if job.Status != "created" {
		t.Errorf("Expected status 'created', got %q", job.Status)
	}
}

func TestJobHandler_GetJob_NotFound(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	req := httptest.NewRequest("GET", "/jobs/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.GetJob(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestJobHandler_CancelJob(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	// First create a job
	reqBody := `{
		"name": "Test Job",
		"sourceUrl": "http://source/file.txt",
		"destination": "target-server",
		"destinationUrl": "http://target/receive"
	}`

	createReq := httptest.NewRequest("POST", "/jobs", strings.NewReader(reqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	handler.CreateJob(createW, createReq)
	location := createW.Result().Header.Get("Location")
	jobID := strings.TrimPrefix(location, "http://localhost:8080/jobs/")

	// Test canceling the job
	cancelReq := httptest.NewRequest("DELETE", "/jobs/"+jobID, nil)
	cancelW := httptest.NewRecorder()
	handler.CancelJob(cancelW, cancelReq)

	resp := cancelW.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Verify the job was cancelled
	getReq := httptest.NewRequest("GET", "/jobs/"+jobID, nil)
	getW := httptest.NewRecorder()
	handler.GetJob(getW, getReq)

	var job models.Job
	if err := json.NewDecoder(getW.Result().Body).Decode(&job); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if job.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got %q", job.Status)
	}
}

func TestJobHandler_CancelJob_NotFound(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	req := httptest.NewRequest("DELETE", "/jobs/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.CancelJob(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestJobHandler_CancelJob_MethodNotAllowed(t *testing.T) {
	service := service.NewJobService()
	handler := NewJobHandler(service, "http://localhost:8080")

	req := httptest.NewRequest("POST", "/jobs/123", nil)
	w := httptest.NewRecorder()

	handler.CancelJob(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
