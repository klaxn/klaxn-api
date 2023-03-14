package data

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Alert struct {
	Item
	UniqueIdentifier string                                `json:"uniq_id" gorm:"index"`
	Status           string                                `json:"status"`
	StartsAt         time.Time                             `json:"starts_at"`
	EndsAt           time.Time                             `json:"ends_at,omitempty"`
	Title            string                                `json:"title"`
	Description      string                                `json:"description"`
	UrlMoreInfo      string                                `json:"url_more_info,omitempty"`
	Labels           datatypes.JSONType[map[string]string] `json:"labels"`
	ServiceID        uint                                  `json:"service_id" gorm:"index"`
}

func (s *Alert) Validate(m *Manager) error {
	if len(s.UniqueIdentifier) == 0 {
		return errors.New("please set a unique identifier for the alert")
	}

	if s.ServiceID == 0 {
		return errors.New("please set an service id in the `service_id` field")
	}

	_, err := m.GetEscalation(s.ServiceID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("no service exists for id %d", s.ServiceID)
	}

	return err
}

func (s *Alert) SafeToDelete(*Manager) error {
	return errors.New("alerts cannot be deleted")
}

func MapToJSONType(in map[string]string) datatypes.JSONType[map[string]string] {
	return datatypes.JSONType[map[string]string]{
		Data: in,
	}
}
