package repository

import (
	"time"

	"github.com/transfer-service/internal/models"
)

type MockJobRepository struct {
	jobs   map[uint]*models.Job
	nextID uint
}

type MockJobRequestRepository struct {
	requests map[uint]*models.JobRequest
	nextID   uint
}

func NewMockJobRepository() JobRepository {
	return &MockJobRepository{
		jobs:   make(map[uint]*models.Job),
		nextID: 1,
	}
}

func NewMockJobRequestRepository() JobRequestRepository {
	return &MockJobRequestRepository{
		requests: make(map[uint]*models.JobRequest),
		nextID:   1,
	}
}

func (r *MockJobRepository) Create(job *models.Job) error {
	job.ID = r.nextID
	job.CreatedAt = time.Now()
	r.jobs[job.ID] = job
	r.nextID++
	return nil
}

func (r *MockJobRepository) GetByID(id uint) (*models.Job, error) {
	job, exists := r.jobs[id]
	if !exists {
		return nil, ErrRecordNotFound
	}
	return job, nil
}

func (r *MockJobRepository) Update(job *models.Job) error {
	if _, exists := r.jobs[job.ID]; !exists {
		return ErrRecordNotFound
	}
	r.jobs[job.ID] = job
	return nil
}

func (r *MockJobRepository) Delete(id uint) error {
	if _, exists := r.jobs[id]; !exists {
		return ErrRecordNotFound
	}
	delete(r.jobs, id)
	return nil
}

func (r *MockJobRepository) List() ([]*models.Job, error) {
	jobs := make([]*models.Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *MockJobRepository) Migrate() error {
	return nil
}

func (r *MockJobRequestRepository) Create(request *models.JobRequest) error {
	request.ID = r.nextID
	request.CreatedAt = time.Now()
	r.requests[request.ID] = request
	r.nextID++
	return nil
}

func (r *MockJobRequestRepository) GetByID(id uint) (*models.JobRequest, error) {
	request, exists := r.requests[id]
	if !exists {
		return nil, ErrRecordNotFound
	}
	return request, nil
}

func (r *MockJobRequestRepository) Update(request *models.JobRequest) error {
	if _, exists := r.requests[request.ID]; !exists {
		return ErrRecordNotFound
	}
	r.requests[request.ID] = request
	return nil
}

func (r *MockJobRequestRepository) Delete(id uint) error {
	if _, exists := r.requests[id]; !exists {
		return ErrRecordNotFound
	}
	delete(r.requests, id)
	return nil
}

func (r *MockJobRequestRepository) List() ([]*models.JobRequest, error) {
	requests := make([]*models.JobRequest, 0, len(r.requests))
	for _, request := range r.requests {
		requests = append(requests, request)
	}
	return requests, nil
}

func (r *MockJobRequestRepository) Migrate() error {
	return nil
}
