package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/klaxn/klaxn-api/pkg/model/team"
)

// GetTeams godoc
// @Summary Get all teams
// @Description Get all teams
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {array} team.Team
// @Failure 500 {object} data.Error
// @Router /teams [get]
func (r *Router) GetTeams(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetTeams")
	defer span.End()

	teams, err := r.db.GetTeams()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}
	r.SendJsonResponse(c, http.StatusOK, team.FromDataSlice(teams), span)
}

// GetTeam godoc
// @Summary Get a team
// @Description Get a team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} team.Team
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /teams/{id} [get]
// @Param id path int true "id"
func (r *Router) GetTeam(c *gin.Context) {
	_, span := r.tracer.Start(c, "GetTeam")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	t, err := r.db.GetTeam(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}
	r.SendJsonResponse(c, http.StatusOK, team.FromData(t), span)
}

// CreateTeam godoc
// @Summary Create a team
// @Description Create a team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} team.Team
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /teams [post]
// @Param team body team.Team true "team"
func (r *Router) CreateTeam(c *gin.Context) {
	_, span := r.tracer.Start(c, "CreateTeam")
	defer span.End()

	var json team.Team
	if err := c.ShouldBindJSON(&json); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	t := json.ToData()

	if err := t.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	newTeam, err := r.db.UpdateTeam(t)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, team.FromData(newTeam), span)
}

// UpdateTeam godoc
// @Summary Update a team
// @Description Update a team
// @Tags teams
// @Accept json
// @Produce json
// @Success 200 {object} team.Team
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /services/{id} [put]
// @Param team body team.Team true "team"
// @Param id path int true "id"
func (r *Router) UpdateTeam(c *gin.Context) {
	_, span := r.tracer.Start(c, "UpdateTeam")
	defer span.End()

	var json team.Team
	if err := c.ShouldBindJSON(&json); err != nil {
		r.SendErr(c, http.StatusBadRequest, errors.Wrap(err, "could not parse json"), span)
		return
	}

	t := json.ToData()

	if err := t.Validate(r.db); err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	t.ID = uint(id)
	s, err := r.db.UpdateTeam(t)
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, team.FromData(s), span)
	return
}

// DeleteTeam godoc
// @Summary Delete a team
// @Description Delete a team
// @Tags teams
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /teams/{id} [delete]
// @Param id path int true "id"
func (r *Router) DeleteTeam(c *gin.Context) {
	_, span := r.tracer.Start(c, "DeleteTeam")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	err = r.db.DeleteTeam(uint(id))
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendNullResponse(c, http.StatusNoContent, span)
	return
}
