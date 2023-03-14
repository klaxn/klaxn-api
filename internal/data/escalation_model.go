package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

const (
	ResponderTypeUser     = "user"
	ResponderTypeSchedule = "schedule"
)

type Escalation struct {
	Item
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Layers      *datatypes.JSONType[[]Layer] `json:"layers,omitempty"`
}

type Layer struct {
	Tier               uint   `json:"tier"`
	ResponderType      string `json:"responder_type"`
	ResponderReference string `json:"responder_reference"`
}

func (e *Escalation) SafeToDelete(m *Manager) error {
	services, err := m.GetServicesOnEscalation(e.ID)
	if err != nil {
		return err
	}

	var errStrs []string

	for _, service := range services {
		errStrs = append(errStrs, fmt.Sprintf("%s (id: %d)", service.Name, service.ID))
	}

	if len(errStrs) > 0 {
		return fmt.Errorf("the following services are asigned to the escalation '%s' therefore, the escalation cannot be deleted: [ %s ]", e.Name, strings.Join(errStrs, ","))
	}
	return nil
}

func (e *Escalation) GetLayers() ([]Layer, error) {
	var layers []Layer
	bytes, err := e.Layers.MarshalJSON()
	if err != nil {
		return layers, err
	}
	if err := json.Unmarshal(bytes, &layers); err != nil {
		return nil, err
	}

	return layers, nil
}

func (l *Layer) Validate(*Manager) error {
	if l.ResponderType != ResponderTypeUser && l.ResponderType != ResponderTypeSchedule {
		return fmt.Errorf("responder type must be one of: [%s, %s]", ResponderTypeUser, ResponderTypeSchedule)
	}

	if len(l.ResponderReference) == 0 {
		return errors.New("please set a responder reference")
	}

	return nil
}

func (l *Layer) SafeToDelete(*Manager) error {
	return errors.New("implement me")
}

func (e *Escalation) Validate(*Manager) error {
	if len(e.Name) == 0 {
		return errors.New("please set a name")
	}

	if len(e.Description) == 0 {
		return errors.New("please set a description")
	}

	layers, err := e.GetLayers()
	if err != nil {
		return err
	}

	if len(layers) == 0 {
		return errors.New("please set at least one layer")
	}

	for i, layer := range layers {
		err := layer.Validate(nil)
		if err != nil {
			return errors.Wrapf(err, "error when validaing layer %d", i)
		}
	}

	return nil
}
