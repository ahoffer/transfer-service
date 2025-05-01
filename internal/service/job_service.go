package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/repository"
)

type JobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}

func (s *JobService) CreateJob(req *models.JobRequest) (*models.Job, error) {
	job := &models.Job{
		Name:           req.Name,
		SourceUrl:      req.SourceUrl,
		Destination:    req.Destination,
		DestinationUrl: req.DestinationUrl,
		Status:         "created",
	}

	if err := s.repo.Create(job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return job, nil
}

func (s *JobService) GetJob(id uint) (*models.Job, error) {
	job, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

func (s *JobService) CancelJob(id uint) error {
	job, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.Status = "cancelled"
	if err := s.repo.Update(job); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

func (s *JobService) ListJobs() ([]*models.Job, error) {
	return s.repo.List()
}

func generateJobID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// If we can't generate random bytes, use timestamp as fallback
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
