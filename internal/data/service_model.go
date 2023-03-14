package data

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	Item
	Name         string `json:"name" gorm:"index,unique"`
	TeamID       uint   `json:"team_id" gorm:"index"`
	Description  string `json:"description"`
	Link         string `json:"link,omitempty"`
	EscalationID uint   `json:"escalation_id" gorm:"index"`
}

func (s *Service) Validate(m *Manager) error {
	if len(s.Name) == 0 {
		return errors.New("please set a service name in the 'name' field")
	}

	if s.EscalationID == 0 {
		return errors.New("please set an escalation id in the `escalation_id` field")
	}

	_, err := m.GetEscalation(s.EscalationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("no escalation exists for id %d", s.EscalationID)
	}
	if err != nil {
		return err
	}

	if s.TeamID == 0 {
		return errors.New("please set an team id in the `team_id` field")
	}

	_, err = m.GetTeam(s.TeamID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("no team exists for id %d", s.TeamID)
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SafeToDelete(*Manager) error {
	return nil
}
