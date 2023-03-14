package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/klaxn/klaxn-api/pkg/model/escalation"
)

// @BasePath /api/v1

// GetEscalations godoc
// @Summary Get all escalations
// @Schemes
// @Description Get all escalations
// @Tags escalations
// @Accept json
// @Produce json
// @Success 200 {array} escalation.Escalation
// @Failure 500 {object} data.Error
// @Router /escalations [get]
func (r *Router) GetEscalations(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetEscalations")
	defer span.End()

	escalations, err := r.db.GetEscalations()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	slice, err := escalation.FromDataSlice(escalations)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, slice, span)
}

// GetEscalation godoc
// @Summary Get an escalation
// @Description Get an escalation
// @Tags escalations
// @Accept json
// @Produce json
// @Success 200 {object} escalation.Escalation
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /escalations/{id} [get]
// @Param id path int true "id"
func (r *Router) GetEscalation(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetEscalation")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	de, err := r.db.GetEscalation(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	e, err := escalation.FromData(de)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	c.JSON(http.StatusOK, e)
}

// CreateEscalation godoc
// @Summary Create an escalation
// @Description Create an escalation
// @Tags escalations
// @Accept json
// @Produce json
// @Success 200 {object} escalation.Escalation
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /escalations [post]
// @Param escalation body escalation.Escalation true "escalation"
func (r *Router) CreateEscalation(c *gin.Context) {
	_, span := r.tracer.Start(c, "CreateEscalation")
	defer span.End()

	var e escalation.Escalation
	if err := c.ShouldBindJSON(&e); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	json, err := e.ToData()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	if err := json.Validate(nil); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	newData, err := r.db.UpdateEscalation(json)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	newEscalation, err := escalation.FromData(newData)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, newEscalation, span)
	return
}

// UpdateEscalation godoc
// @Summary Update an existing escalation
// @Description Create an existing escalation
// @Tags escalations
// @Accept json
// @Produce json
// @Success 200 {object} escalation.Escalation
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /escalations/{id} [put]
// @Param escalation body escalation.Escalation true "escalation"
// @Param id path int true "id"
func (r *Router) UpdateEscalation(c *gin.Context) {
	_, span := r.tracer.Start(c, "UpdateEscalation")
	defer span.End()

	var e escalation.Escalation
	if err := c.ShouldBindJSON(&e); err != nil {
		r.SendErr(c, http.StatusBadRequest, errors.Wrap(err, "could not parse json"), span)
		return
	}

	json, err := e.ToData()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	if err := json.Validate(nil); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	newData, err := r.db.UpdateEscalation(json)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	updatedEscalation, err := escalation.FromData(newData)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, updatedEscalation, span)
	return
}

// DeleteEscalation godoc
// @Summary Delete an existing escalation
// @Description Delete an existing escalation
// @Tags escalations
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /escalations/{id} [delete]
// @Param id path int true "id"
func (r *Router) DeleteEscalation(c *gin.Context) {
	_, span := r.tracer.Start(c, "DeleteEscalation")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	err = r.db.DeleteEscalation(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendNullResponse(c, http.StatusNoContent, span)
	return
}
