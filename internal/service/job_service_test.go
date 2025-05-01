package service

import (
	"testing"

	"github.com/transfer-service/internal/models"
	"github.com/transfer-service/internal/repository"
)

func TestJobService_CreateJob(t *testing.T) {
	repo := repository.NewMockJobRepository()
	service := NewJobService(repo)

	req := &models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	job, err := service.CreateJob(req)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	// Verify job fields
	if job.Name != req.Name {
		t.Errorf("Expected job name %q, got %q", req.Name, job.Name)
	}
	if job.SourceUrl != req.SourceUrl {
		t.Errorf("Expected source URL %q, got %q", req.SourceUrl, job.SourceUrl)
	}
	if job.Destination != req.Destination {
		t.Errorf("Expected destination %q, got %q", req.Destination, job.Destination)
	}
	if job.DestinationUrl != req.DestinationUrl {
		t.Errorf("Expected destination URL %q, got %q", req.DestinationUrl, job.DestinationUrl)
	}
	if job.Status != "created" {
		t.Errorf("Expected status 'created', got %q", job.Status)
	}
	if job.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestJobService_GetJob(t *testing.T) {
	repo := repository.NewMockJobRepository()
	service := NewJobService(repo)

	// Create a job first
	req := &models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	createdJob, err := service.CreateJob(req)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	// Test getting the created job
	job, err := service.GetJob(createdJob.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}

	// Verify the job matches what we created
	if job.ID != createdJob.ID {
		t.Errorf("Expected job ID %d, got %d", createdJob.ID, job.ID)
	}
	if job.Name != req.Name {
		t.Errorf("Expected job name %q, got %q", req.Name, job.Name)
	}
	if job.SourceUrl != req.SourceUrl {
		t.Errorf("Expected source URL %q, got %q", req.SourceUrl, job.SourceUrl)
	}
	if job.Destination != req.Destination {
		t.Errorf("Expected destination %q, got %q", req.Destination, job.Destination)
	}
	if job.DestinationUrl != req.DestinationUrl {
		t.Errorf("Expected destination URL %q, got %q", req.DestinationUrl, job.DestinationUrl)
	}

	// Test getting non-existent job
	_, err = service.GetJob(999)
	if err == nil {
		t.Error("Expected error when getting non-existent job")
	}
}

func TestJobService_CancelJob(t *testing.T) {
	repo := repository.NewMockJobRepository()
	service := NewJobService(repo)

	// Create a job first
	req := &models.JobRequest{
		Name:           "Test Job",
		SourceUrl:      "http://source/file.txt",
		Destination:    "target-server",
		DestinationUrl: "http://target/receive",
	}

	createdJob, err := service.CreateJob(req)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	// Test canceling the job
	err = service.CancelJob(createdJob.ID)
	if err != nil {
		t.Fatalf("CancelJob failed: %v", err)
	}

	// Verify the job status was updated
	job, err := service.GetJob(createdJob.ID)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if job.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got %q", job.Status)
	}

	// Test canceling non-existent job
	err = service.CancelJob(999)
	if err == nil {
		t.Error("Expected error when canceling non-existent job")
	}
}
