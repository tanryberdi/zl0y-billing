package handlers

import (
	"net/http"

	"zl0y-billing/internal/models"
	"zl0y-billing/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) PurchaseReport(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	reportID := c.Param("report_id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Report ID is required",
		})
		return
	}

	err := h.reportService.PurchaseReport(userID.(int), reportID)
	if err != nil {
		switch err.Error() {
		case "report not found":
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Report not found",
			})
		case "report already purchased":
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: "Report already purchased",
			})
		case "insufficient balance":
			c.JSON(http.StatusPaymentRequired, models.ErrorResponse{
				Error: "Insufficient balance",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to purchase report",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Report purchased successfully",
	})
}
