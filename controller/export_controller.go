package controller

import (
    "eventix/service"
    "encoding/csv"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type ExportController struct {
    reportService service.ReportService
}

func NewExportController(reportService service.ReportService) *ExportController {
    return &ExportController{reportService: reportService}
}

func (ctrl *ExportController) ExportSummaryReport(c *gin.Context) {
    report, err := ctrl.reportService.GetSummaryReport()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate summary report"})
        return
    }

    // Membuat CSV
    c.Writer.Header().Set("Content-Type", "text/csv")
    c.Writer.Header().Set("Content-Disposition", "attachment;filename=summary_report.csv")
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    // Header CSV
    writer.Write([]string{"Total Tickets Sold", "Total Revenue"})

    // Data CSV
    writer.Write([]string{
        strconv.FormatInt(report["total_tickets_sold"].(int64), 10),
        strconv.FormatFloat(report["total_revenue"].(float64), 'f', 2, 64),
    })
}

func (ctrl *ExportController) ExportEventReport(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID"})
        return
    }

    report, err := ctrl.reportService.GetEventReport(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate event report"})
        return
    }

    // Membuat CSV
    c.Writer.Header().Set("Content-Type", "text/csv")
    c.Writer.Header().Set("Content-Disposition", "attachment;filename=event_report.csv")
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()

    // Header CSV
    writer.Write([]string{"Total Tickets Sold", "Total Revenue"})

    // Data CSV
    writer.Write([]string{
        strconv.FormatInt(report["total_tickets_sold"].(int64), 10),
        strconv.FormatFloat(report["total_revenue"].(float64), 'f', 2, 64),
    })
}
