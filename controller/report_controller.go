package controller

import (
	"eventix/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
    reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
    return &ReportController{
        reportService: reportService,
    }
}

func (ctrl *ReportController) GetSalesReport(c *gin.Context) {
    report, err := ctrl.reportService.GenerateSalesReport()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "success", "data": report})
}

func (ctrl *ReportController) GetSummaryReport(c *gin.Context) {
    report, err := ctrl.reportService.GetSummaryReport()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate summary report", "data": nil})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Summary report generated successfully", "data": report})
}

func (ctrl *ReportController) GetEventReport(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID", "data": nil})
        return
    }

    report, err := ctrl.reportService.GetEventReport(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate event report", "data": nil})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event report generated successfully", "data": report})
}
