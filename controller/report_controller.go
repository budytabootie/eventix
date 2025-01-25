package controller

import (
	"eventix/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service service.ReportService
}

func NewReportController(service service.ReportService) *ReportController {
	return &ReportController{
		service: service,
	}
}

func (ctrl *ReportController) GetEventReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID", "data": nil})
		return
	}

	report, err := ctrl.service.GetEventReport(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event report retrieved successfully", "data": report})
}

func (ctrl *ReportController) GetSummaryReport(c *gin.Context) {
	report, err := ctrl.service.GetSummaryReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Summary report retrieved successfully", "data": report})
}
