package handlers

import (
	"net/http"

	"zl0y-billing/internal/models"
	"zl0y-billing/internal/repository"

	"github.com/gin-gonic/gin"
)

type MockHandler struct {
	reportRepo *repository.ReportRepository
}

func NewMockHandler(reportRepo *repository.ReportRepository) *MockHandler {
	return &MockHandler{
		reportRepo: reportRepo,
	}
}

func (h *MockHandler) CreateReport(c *gin.Context) {
	var req models.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	report, err := h.reportRepo.CreateReport(req.ClientGeneratedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create report",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Report created successfully",
		"report_id": report.ReportID,
		"report":    report,
	})
}
