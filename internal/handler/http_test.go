package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/repository"
	"github.com/transfer-service/internal/service"
)

type mockJobRepository struct {
	jobs map[uint]*models.Job
}

func (m *mockJobRepository) Create(job *models.Job) error {
	job.ID = uint(len(m.jobs) + 1)
	m.jobs[job.ID] = job
	return nil
}

func (m *mockJobRepository) GetByID(id uint) (*models.Job, error) {
	if job, exists := m.jobs[id]; exists {
		return job, nil
	}
	return nil, repository.ErrRecordNotFound
}

func (m *mockJobRepository) Update(job *models.Job) error {
	if _, exists := m.jobs[job.ID]; !exists {
		return repository.ErrRecordNotFound
	}
	m.jobs[job.ID] = job
	return nil
}

func (m *mockJobRepository) Delete(id uint) error {
	if _, exists := m.jobs[id]; !exists {
		return repository.ErrRecordNotFound
	}
	delete(m.jobs, id)
	return nil
}

func (m *mockJobRepository) List() ([]*models.Job, error) {
	jobs := make([]*models.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (m *mockJobRepository) Migrate() error {
	return nil
}

type mockJobRequestRepository struct {
	requests map[uint]*models.JobRequest
}

func (m *mockJobRequestRepository) Create(req *models.JobRequest) error {
	req.ID = uint(len(m.requests) + 1)
	m.requests[req.ID] = req
	return nil
}

func (m *mockJobRequestRepository) GetByID(id uint) (*models.JobRequest, error) {
	if req, exists := m.requests[id]; exists {
		return req, nil
	}
	return nil, repository.ErrRecordNotFound
}

func (m *mockJobRequestRepository) List() ([]*models.JobRequest, error) {
	requests := make([]*models.JobRequest, 0, len(m.requests))
	for _, req := range m.requests {
		requests = append(requests, req)
	}
	return requests, nil
}

func (m *mockJobRequestRepository) Migrate() error {
	return nil
}

func TestJobEndpoints(t *testing.T) {
	// Initialize mock repositories
	jobRepo := &mockJobRepository{jobs: make(map[uint]*models.Job)}
	requestRepo := &mockJobRequestRepository{requests: make(map[uint]*models.JobRequest)}

	// Initialize the job service and handler
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	// Test Create Job
	t.Run("Create Job", func(t *testing.T) {
		jobReq := models.JobRequest{
			Name:           "Test Job",
			SourceUrl:      "http://example.com/source",
			Destination:    "s3://bucket/path",
			DestinationUrl: "http://example.com/destination",
		}

		reqBody, err := json.Marshal(jobReq)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "/jobs/")
	})

	// Test Get Job
	t.Run("Get Job", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/jobs/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var job models.Job
		err := json.NewDecoder(w.Body).Decode(&job)
		assert.NoError(t, err)
		assert.Equal(t, "Test Job", job.Name)
	})

	// Test Cancel Job
	t.Run("Cancel Job", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/jobs/1/cancel", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestJobValidation(t *testing.T) {
	// Initialize mock repositories
	jobRepo := &mockJobRepository{jobs: make(map[uint]*models.Job)}
	requestRepo := &mockJobRequestRepository{requests: make(map[uint]*models.JobRequest)}

	// Initialize the job service and handler
	jobService := service.NewJobService(jobRepo, requestRepo)
	handler := NewJobHandler(jobService, "http://localhost:8080")
	router := setupTestRouter(handler)

	tests := []struct {
		name        string
		jobReq      models.JobRequest
		expectedErr bool
	}{
		{
			name: "Valid Job Request",
			jobReq: models.JobRequest{
				Name:           "Valid Job",
				SourceUrl:      "http://example.com/source",
				Destination:    "s3://bucket/path",
				DestinationUrl: "http://example.com/destination",
			},
			expectedErr: false,
		},
		{
			name: "Invalid Job Request - Missing Name",
			jobReq: models.JobRequest{
				SourceUrl:      "http://example.com/source",
				Destination:    "s3://bucket/path",
				DestinationUrl: "http://example.com/destination",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.jobReq)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectedErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
			}
		})
	}
}
