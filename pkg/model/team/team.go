package team

import (
	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/pkg/model"
)

type Team struct {
	model.Item
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FromData(d *data.Team) *Team {
	e := &Team{
		Item:        model.Item{ID: d.ID},
		Name:        d.Name,
		Description: d.Description,
	}

	return e
}

func (e *Team) ToData() *data.Team {
	d := &data.Team{
		Item:        data.Item{ID: e.ID},
		Name:        e.Name,
		Description: e.Description,
	}

	return d
}

func FromDataSlice(in []*data.Team) []*Team {
	var out []*Team

	for _, service := range in {
		fromData := FromData(service)
		out = append(out, fromData)
	}

	return out
}
