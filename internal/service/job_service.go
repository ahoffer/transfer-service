package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/transfer-service/internal/models"
)

type JobService struct {
	jobs map[string]*models.Job
}

func NewJobService() *JobService {
	return &JobService{
		jobs: make(map[string]*models.Job),
	}
}

func (s *JobService) CreateJob(req *models.JobRequest) (*models.Job, error) {
	id := generateJobID()

	job := &models.Job{
		ID:         id,
		Status:     "created",
		CreatedAt:  time.Now(),
		JobRequest: *req,
	}

	s.jobs[id] = job
	return job, nil
}

func (s *JobService) GetJob(id string) (*models.Job, error) {
	job, exists := s.jobs[id]
	if !exists {
		return nil, errors.New("job not found")
	}
	return job, nil
}

func (s *JobService) CancelJob(id string) error {
	job, exists := s.jobs[id]
	if !exists {
		return errors.New("job not found")
	}
	job.Status = "cancelled"
	return nil
}

func generateJobID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// If we can't generate random bytes, use timestamp as fallback
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
