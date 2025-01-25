package controller

import (
	"eventix/entity"
	"eventix/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	service service.EventService
}

func NewEventController(eventService service.EventService) *EventController {
	return &EventController{
		service: eventService,
	}
}


// GetAllEvents godoc
// @Summary Get all events
// @Description Retrieve a list of all events (accessible by Users)
// @Tags Events
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events [get]
func (ctrl *EventController) GetAllEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	events, err := ctrl.service.GetAllEvents(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve events", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Events retrieved successfully", "data": events})
}

func (ctrl *EventController) GetEventByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID", "data": nil})
		return
	}
	event, err := ctrl.service.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Event not found", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event retrieved successfully", "data": event})
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Admin can create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Param event body entity.Event true "Event details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/events [post]
func (ctrl *EventController) CreateEvent(c *gin.Context) {
	var event entity.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event data", "data": nil})
		return
	}
	createdEvent, err := ctrl.service.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Event created successfully", "data": createdEvent})
}

func (ctrl *EventController) UpdateEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID", "data": nil})
		return
	}

	var event entity.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event data", "data": nil})
		return
	}
	event.ID = uint(id)

	updatedEvent, err := ctrl.service.UpdateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event updated successfully", "data": updatedEvent})
}

func (ctrl *EventController) DeleteEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid event ID", "data": nil})
		return
	}
	if err := ctrl.service.DeleteEvent(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Event deleted successfully", "data": nil})
}
