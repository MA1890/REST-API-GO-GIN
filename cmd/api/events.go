package main

import (
	"REST-API-GO-GIN/internal/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := app.GetUserFromContext(c)
	user.Id = event.OwnerId
	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}
	event, err := app.models.Events.Get(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve event"})
		return
	}
	c.JSON(http.StatusOK, event)
}

// GetEvents returns all events
//
//	@Summary		Returns all events
//	@Description	Returns all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]database.Event
//	@Router			/api/v1/events [get]

func (app *application) getAllEvents(c *gin.Context) {

	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve all events"})
		return
	}
	c.JSON(http.StatusOK, events)

}

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return

	}
	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update"})
	}
	updateEvent := &database.Event{
		Id: id,
	}
	if err := c.ShouldBindJSON(updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = app.models.Events.Update(updateEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Update event"})
		return
	}
	c.JSON(http.StatusOK, updateEvent)

}

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return

	}
	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Delete event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AddAttendeeToEvent adds an attendee to an event
// @Summary		Adds an attendee to an event
// @Description	Adds an attendee to an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		201		{object}	database.Attendee
// @Router			/api/v1/events/{id}/attendees/{userId} [post]
// @Security		BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve user to add event"})
		return

	}
	if userToAdd.Id != event.OwnerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add event"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(eventId, userToAdd.Id)
	if err != nil {
		log.Printf("Error retrieving attendee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve user to add event"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
		return
	}
	attendee := database.Attendee{
		EventId: eventId,
		UserId:  userId,
	}
	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Add user to event"})
		return
	}
	c.JSON(http.StatusCreated, attendee)

}

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary		Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}
	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary		Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags			attendees
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success		204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth
func (app *application) deleteAttendeeForEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Event ID"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Retrieve event"})
		return
	}
	user := app.GetUserFromContext(c)

	if user.Id != event.OwnerId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	err = app.models.Attendees.Delete(userId, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete User"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary		Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}
	events, err := app.models.Events.GetByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}
