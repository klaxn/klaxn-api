package service

import (
	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/pkg/model"
)

type Service struct {
	model.Item
	Name         string `json:"name"`
	TeamID       uint   `json:"team_id"`
	Description  string `json:"description"`
	Link         string `json:"link,omitempty"`
	EscalationID uint   `json:"escalation_id"`
}

func FromData(d *data.Service) *Service {
	e := &Service{
		Item:         model.Item{ID: d.ID},
		Name:         d.Name,
		Description:  d.Description,
		Link:         d.Link,
		EscalationID: d.EscalationID,
		TeamID:       d.TeamID,
	}

	return e
}

func (e *Service) ToData() *data.Service {
	d := &data.Service{
		Item:         data.Item{ID: e.ID},
		Name:         e.Name,
		Description:  e.Description,
		Link:         e.Link,
		EscalationID: e.EscalationID,
		TeamID:       e.TeamID,
	}

	return d
}

func FromDataSlice(in []*data.Service) []*Service {
	var out []*Service

	for _, service := range in {
		fromData := FromData(service)
		out = append(out, fromData)
	}

	return out
}
