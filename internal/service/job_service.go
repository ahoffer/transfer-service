package service

import (
	"errors"
	"fmt"

	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/repository"
)

type JobService struct {
	jobRepo     repository.JobRepository
	requestRepo repository.JobRequestRepository
}

func NewJobService(jobRepo repository.JobRepository, requestRepo repository.JobRequestRepository) *JobService {
	return &JobService{
		jobRepo:     jobRepo,
		requestRepo: requestRepo,
	}
}

func (s *JobService) CreateJob(req *models.JobRequest) (*models.Job, error) {
	// First, persist the request
	if err := s.requestRepo.Create(req); err != nil {
		return nil, fmt.Errorf("failed to create job request: %w", err)
	}

	// Then create the job
	job := &models.Job{
		Name:           req.Name,
		SourceUrl:      req.SourceUrl,
		Destination:    req.Destination,
		DestinationUrl: req.DestinationUrl,
		Status:         "created",
	}

	if err := s.jobRepo.Create(job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return job, nil
}

func (s *JobService) GetJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, errors.New("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

func (s *JobService) CancelJob(id uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return errors.New("job not found")
		}
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.Status = "cancelled"
	if err := s.jobRepo.Update(job); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

func (s *JobService) ListJobs() ([]*models.Job, error) {
	return s.jobRepo.List()
}
