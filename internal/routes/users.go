package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/klaxn/klaxn-api/pkg/model/user"
)

// GetUsers godoc
// @Summary Get all users
// @Schemes
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} user.User
// @Failure 500 {object} data.Error
// @Router /users [get]
func (r *Router) GetUsers(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetUsers")
	defer span.End()

	u, err := r.db.GetUsers()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, user.FromDataSlice(u), span)
}

// GetUser godoc
// @Summary Get a user
// @Description Get a user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.User
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /users/{id} [get]
// @Param id path int true "id"
func (r *Router) GetUser(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetUser")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	u, err := r.db.GetUser(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, user.FromData(u), span)
}

// CreateUser godoc
// @Summary Create a user
// @Description Create a user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.User
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /users [post]
// @Param user body user.User true "user"
func (r *Router) CreateUser(c *gin.Context) {
	_, span := r.tracer.Start(c, "CreateUser")
	defer span.End()

	var json user.User
	if err := c.ShouldBindJSON(&json); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	u := json.ToData()

	if err := u.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	createdUser, err := r.db.CreateUser(u)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, user.FromData(createdUser), span)
	return
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Create an existing user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} user.User
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /users/{id} [put]
// @Param user body user.User true "user"
// @Param id path int true "id"
func (r *Router) UpdateUser(c *gin.Context) {
	_, span := r.tracer.Start(c, "UpdateUser")
	defer span.End()

	var json user.User
	if err := c.ShouldBindJSON(&json); err != nil {
		r.SendErr(c, http.StatusBadRequest, errors.Wrap(err, "could not parse json"), span)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	u := json.ToData()

	u.ID = uint(id)

	if err := u.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	updatedUser, err := r.db.UpdateUser(u)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, user.FromData(updatedUser), span)
	return
}

// DeleteUser godoc
// @Summary Delete an existing user
// @Description Delete an existing user
// @Tags users
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /users/{id} [delete]
// @Param id path int true "id"
func (r *Router) DeleteUser(c *gin.Context) {
	_, span := r.tracer.Start(c, "DeleteUser")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	err = r.db.DeleteUser(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendNullResponse(c, http.StatusNoContent, span)
	return
}
