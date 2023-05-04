package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/klaxn/klaxn-api/pkg/model/service"
)

// GetServices godoc
//
//	@Summary		Get all services
//	@Description	Get all services
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		service.Service
//	@Failure		500	{object}	data.Error
//	@Router			/services [get]
func (r *Router) GetServices(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetServices")
	defer span.End()

	services, err := r.db.GetServices()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, service.FromDataSlice(services), span)
}

// GetService godoc
//
//	@Summary		Get a service
//	@Description	Get a service
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	service.Service
//	@Failure		400	{object}	data.Error
//	@Failure		500	{object}	data.Error
//	@Router			/services/{id} [get]
//	@Param			id	path	int	true	"id"
func (r *Router) GetService(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetService")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	s, err := r.db.GetService(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, service.FromData(s), span)
}

// CreateService godoc
//
//	@Summary		Create a service
//	@Description	Create a service
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	service.Service
//	@Failure		400	{object}	data.Error
//	@Failure		500	{object}	data.Error
//	@Router			/services [post]
//	@Param			service	body	service.Service	true	"service"
func (r *Router) CreateService(c *gin.Context) {
	_, span := r.tracer.Start(c, "CreateService")
	defer span.End()

	var json service.Service
	if err := c.ShouldBindJSON(&json); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	toData := json.ToData()
	if err := toData.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	s, err := r.db.UpdateService(toData)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, service.FromData(s), span)
}

// UpdateService godoc
//
//	@Summary		Update a service
//	@Description	Update a service
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	service.Service
//	@Failure		400	{object}	data.Error
//	@Failure		500	{object}	data.Error
//	@Router			/services/{id} [put]
//	@Param			service	body	service.Service	true	"service"
//	@Param			id		path	int				true	"id"
func (r *Router) UpdateService(c *gin.Context) {
	_, span := r.tracer.Start(c, "UpdateService")
	defer span.End()

	var s service.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		r.SendErr(c, http.StatusBadRequest, errors.Wrap(err, "could not parse json"), span)
		return
	}

	json := s.ToData()

	if err := json.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	json.ID = uint(id)
	ss, err := r.db.UpdateService(json)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	fromData := service.FromData(ss)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, fromData, span)
	return
}

// DeleteService godoc
//
//	@Summary		Delete a service
//	@Description	Delete a service
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	data.Error
//	@Failure		500	{object}	data.Error
//	@Router			/services/{id} [delete]
//	@Param			id	path	int	true	"id"
func (r *Router) DeleteService(c *gin.Context) {
	_, span := r.tracer.Start(c, "DeleteService")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	err = r.db.DeleteService(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendNullResponse(c, http.StatusNoContent, span)
	return
}
