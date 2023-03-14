package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/klaxn/klaxn-api/internal/data"
)

type GrafanaAlert struct {
	Status       string                 `json:"status"`
	Labels       map[string]string      `json:"labels"`
	Annotations  map[string]string      `json:"annotations"`
	StartsAt     time.Time              `json:"startsAt"`
	EndsAt       time.Time              `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Fingerprint  string                 `json:"fingerprint"`
	SilenceURL   string                 `json:"silenceURL"`
	DashboardURL string                 `json:"dashboardURL"`
	PanelURL     string                 `json:"panelURL"`
	ValueString  string                 `json:"valueString"`
	Values       map[string]interface{} `json:"values"`
}

type GrafanaBody struct {
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	Alerts            []GrafanaAlert    `json:"alerts"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	OrgId             int               `json:"orgId"`
	Title             string            `json:"title"`
	State             string            `json:"state"`
	Message           string            `json:"message"`
}

// GrafanaInbound godoc
// @Summary Send a Grafana alert
// @Description Send a Grafana alert
// @Tags alerts
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} data.Error
// @Failure 500 {object} data.Error
// @Router /alerts/grafana [post]
// @Param alert body GrafanaBody true "aolert"
func (r *Router) GrafanaInbound(c *gin.Context) {
	_, span := r.tracer.Start(c, "GrafanaInbound")
	defer span.End()

	var body GrafanaBody
	err := c.Bind(&body)
	if err != nil {
		r.SendErr(c, http.StatusBadRequest, err, span)
		return
	}

	if body.Message == "" {
		r.SendErr(c, http.StatusBadRequest, errors.New("alert has no message"), span)
		return
	}
	r.logger.Infof("got body: %s", body)

	for _, alert := range body.Alerts {

		serviceName := tryStrings(alert.Annotations["Service"], body.CommonAnnotations["Service"])
		s, err := r.db.GetServiceByName(serviceName)
		if err != nil {
			r.SendNullResponse(c, http.StatusInternalServerError, span)
			return
		}

		if s == nil {
			r.SendJsonResponse(c, http.StatusBadRequest, data.Error{Message: fmt.Sprintf("could not find a service called %s", serviceName)}, span)
			return
		}

		newAlert := data.Alert{
			UniqueIdentifier: alert.Fingerprint,
			Status:           alert.Status,
			StartsAt:         alert.StartsAt,
			EndsAt:           alert.EndsAt,
			Title:            tryStrings(alert.Annotations["summary"], body.CommonAnnotations["summary"]),
			Description:      tryStrings(alert.Annotations["description"], body.CommonAnnotations["description"]),
			UrlMoreInfo:      alert.GeneratorURL,
			Labels:           data.MapToJSONType(alert.Labels),
			ServiceID:        s.ID,
		}

		err = newAlert.Validate(r.db)
		if err != nil {
			r.SendErr(c, http.StatusBadRequest, err, span)
			return
		}
		l := r.logger.WithField("service", serviceName).WithField("fingerprint", alert.Fingerprint)

		l.Info("checking to see if this alert has been seen before")
		existingAlert, err := r.db.GetAlertByUniqueID(alert.Fingerprint)
		if err != nil {
			r.SendErr(c, http.StatusInternalServerError, err, span)
			return
		}

		if existingAlert == nil {
			l.Info("this alert HAS NOT been seen before")
			_, err := r.db.UpdateAlert(&newAlert)
			if err != nil {
				r.SendErr(c, http.StatusInternalServerError, err, span)
				return
			}
		} else {
			l.Info("this alert HAS been seen before")

		}

	}

	for _, sender := range r.ob {
		err = sender.SendMessage("07538798561", fmt.Sprintf("Incoming page: %s", body.Message))
		if err != nil {
			r.SendErr(c, http.StatusInternalServerError, err, span)
			return
		}

	}
	r.SendNullResponse(c, http.StatusNoContent, span)
}

// GetAlerts godoc
// @Summary Get all alerts
// @Description Get all alerts
// @Tags alerts
// @Accept json
// @Produce json
// @Success 200 {array} data.Alert
// @Failure 500 {object} data.Error
// @Router /services [get]
func (r *Router) GetAlerts(c *gin.Context) {
	// TODO: do a proper abstraction
	_, span := r.tracer.Start(c, "GetServices")
	defer span.End()

	alerts, err := r.db.GetAlerts()
	if err != nil {
		r.SendErr(c, http.StatusInternalServerError, err, span)
		return
	}

	r.SendJsonResponse(c, http.StatusOK, alerts, span)
}

func tryStrings(in ...string) string {
	for _, s := range in {
		if s != "" {
			return s
		}
	}

	return ""
}
