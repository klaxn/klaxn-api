package data

func (m *Manager) GetEscalations() ([]*Escalation, error) {
	var escalations []*Escalation
	tx := m.db.Find(&escalations)
	return escalations, tx.Error
}

func (m *Manager) GetEscalation(id uint) (*Escalation, error) {
	var e Escalation

	tx := m.db.First(&e, id)
	return &e, tx.Error
}

func (m *Manager) UpdateEscalation(e *Escalation) (*Escalation, error) {
	tx := m.db.Save(&e)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m.GetEscalation(e.ID)
}

func (m *Manager) DeleteEscalation(id uint) error {
	escalation, err := m.GetEscalation(id)
	if err != nil {
		return err
	}

	if err := escalation.SafeToDelete(m); err != nil {
		return err
	}
	tx := m.db.Delete(escalation)
	return tx.Error
}
