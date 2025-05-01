package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/transfer-service/internal/models"
)

const baseURL = "http://localhost:8080"

func TestIntegration(t *testing.T) {
	// Wait for service to be ready
	waitForService(t)

	// Test Create Job
	t.Run("Create Job", func(t *testing.T) {
		jobReq := models.JobRequest{
			Name:           "Integration Test Job",
			SourceUrl:      "http://example.com/source",
			Destination:    "s3://bucket/integration-test",
			DestinationUrl: "http://example.com/destination",
		}

		reqBody, err := json.Marshal(jobReq)
		assert.NoError(t, err)

		resp, err := http.Post(baseURL+"/jobs", "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Get the job ID from the Location header
		location := resp.Header.Get("Location")
		assert.NotEmpty(t, location)

		// Test Get Job using the location from create
		t.Run("Get Created Job", func(t *testing.T) {
			resp, err := http.Get(baseURL + location)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var job models.Job
			err = json.NewDecoder(resp.Body).Decode(&job)
			assert.NoError(t, err)
			assert.Equal(t, "Integration Test Job", job.Name)
			assert.Equal(t, "created", job.Status)
		})

		// Test Cancel Job
		t.Run("Cancel Job", func(t *testing.T) {
			resp, err := http.Post(baseURL+location+"/cancel", "application/json", nil)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify job is cancelled
			resp, err = http.Get(baseURL + location)
			assert.NoError(t, err)
			defer resp.Body.Close()

			var job models.Job
			err = json.NewDecoder(resp.Body).Decode(&job)
			assert.NoError(t, err)
			assert.Equal(t, "cancelled", job.Status)
		})
	})

	// Test Invalid Job Creation
	t.Run("Create Invalid Job", func(t *testing.T) {
		jobReq := models.JobRequest{
			// Missing required Name field
			SourceUrl:      "http://example.com/source",
			Destination:    "s3://bucket/path",
			DestinationUrl: "http://example.com/destination",
		}

		reqBody, err := json.Marshal(jobReq)
		assert.NoError(t, err)

		resp, err := http.Post(baseURL+"/jobs", "application/json", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// waitForService attempts to connect to the service for up to 30 seconds
func waitForService(t *testing.T) {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("Service did not become ready within %d seconds", maxAttempts)
}
