package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		baseURL:    baseURL,
	}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var req models.JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if req.SourceUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SourceUrl is required"})
		return
	}
	if req.Destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Destination is required"})
		return
	}
	if req.DestinationUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DestinationUrl is required"})
		return
	}

	job, err := h.jobService.CreateJob(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.Header("Location", fmt.Sprintf("/jobs/%d", job.ID))
	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.jobService.GetJob(uint(id))
	if err != nil {
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) CancelJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	if err := h.jobService.CancelJob(uint(id)); err != nil {
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel job"})
		return
	}

	// Get the updated job to return in response
	job, err := h.jobService.GetJob(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated job"})
		return
	}

	c.JSON(http.StatusOK, job)
}
