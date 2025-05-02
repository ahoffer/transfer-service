package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/repository"
	"github.com/transfer-service/internal/service"
)

func setupTestRouter(handler *JobHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/jobs", handler.CreateJob)
	r.GET("/jobs/:id", handler.GetJob)
	r.POST("/jobs/:id/cancel", handler.CancelJob)
	return r
}

func TestJobHandler_CreateJob(t *testing.T) {
	jobRepo := repository.NewMockJobRepository()
	requestRepo := repository.NewMockJobRequestRepository()
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	jobReq := models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	reqBody, _ := json.Marshal(jobReq)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "/jobs/")
}

func TestJobHandler_GetJob(t *testing.T) {
	jobRepo := repository.NewMockJobRepository()
	requestRepo := repository.NewMockJobRequestRepository()
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	// First create a job
	jobReq := models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	createReqBody, _ := json.Marshal(jobReq)
	createReq := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	location := createW.Header().Get("Location")
	jobID := location[len("/jobs/"):]

	// Now test getting the job
	getReq := httptest.NewRequest(http.MethodGet, "/jobs/"+jobID, nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	assert.Equal(t, http.StatusOK, getW.Code)

	var job models.Job
	err := json.NewDecoder(getW.Body).Decode(&job)
	assert.NoError(t, err)
	assert.Equal(t, "Test Job", job.Name)
	assert.Equal(t, "http://source/file.txt", job.SourceUrl)
	assert.Equal(t, "target-server", job.Destination)
	assert.Equal(t, "http://target/receive", job.DestinationUrl)
	assert.Equal(t, "created", job.Status)
}

func TestJobHandler_GetJob_NotFound(t *testing.T) {
	jobRepo := repository.NewMockJobRepository()
	requestRepo := repository.NewMockJobRequestRepository()
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/jobs/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestJobHandler_CancelJob(t *testing.T) {
	jobRepo := repository.NewMockJobRepository()
	requestRepo := repository.NewMockJobRequestRepository()
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	// First create a job
	jobReq := models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	createReqBody, _ := json.Marshal(jobReq)
	createReq := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(createReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	location := createW.Header().Get("Location")
	jobID := location[len("/jobs/"):]

	// Test canceling the job
	cancelReq := httptest.NewRequest(http.MethodPost, "/jobs/"+jobID+"/cancel", nil)
	cancelW := httptest.NewRecorder()
	router.ServeHTTP(cancelW, cancelReq)

	assert.Equal(t, http.StatusOK, cancelW.Code)

	// Verify the job was cancelled
	getReq := httptest.NewRequest(http.MethodGet, "/jobs/"+jobID, nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	var job models.Job
	err := json.NewDecoder(getW.Body).Decode(&job)
	assert.NoError(t, err)
	assert.Equal(t, "cancelled", job.Status)
}

func TestJobHandler_CancelJob_NotFound(t *testing.T) {
	jobRepo := repository.NewMockJobRepository()
	requestRepo := repository.NewMockJobRequestRepository()
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/jobs/999/cancel", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
