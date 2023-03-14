package escalation

import (
	"gorm.io/datatypes"

	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/pkg/model"
)

type Escalation struct {
	model.Item
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Layers      []Layer `json:"layers,omitempty"`
}

type Layer struct {
	Tier               uint   `json:"tier"`
	ResponderType      string `json:"responder_type"`
	ResponderReference string `json:"responder_reference"`
}

func FromData(d *data.Escalation) (*Escalation, error) {
	e := &Escalation{
		Item:        model.Item{ID: d.ID},
		Name:        d.Name,
		Description: d.Description,
		Layers:      []Layer{},
	}

	layers, err := d.GetLayers()
	if err != nil {
		return nil, err
	}

	for _, layer := range layers {
		e.Layers = append(e.Layers, Layer{
			Tier:               layer.Tier,
			ResponderType:      layer.ResponderType,
			ResponderReference: layer.ResponderReference,
		})
	}

	return e, nil
}

func (e *Escalation) ToData() (*data.Escalation, error) {
	var layers []data.Layer

	for _, layer := range e.Layers {
		layers = append(layers, data.Layer{
			Tier:               layer.Tier,
			ResponderType:      layer.ResponderType,
			ResponderReference: layer.ResponderReference,
		})
	}

	d := &data.Escalation{
		Item:        data.Item{ID: e.ID},
		Name:        e.Name,
		Description: e.Description,
		Layers:      &datatypes.JSONType[[]data.Layer]{Data: layers},
	}

	return d, nil
}

func FromDataSlice(in []*data.Escalation) ([]*Escalation, error) {
	var out []*Escalation

	for _, escalation := range in {
		fromData, err := FromData(escalation)
		if err != nil {
			return nil, err
		}
		out = append(out, fromData)
	}

	return out, nil
}
