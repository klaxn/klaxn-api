package data

func (m *Manager) UpdateAlert(a *Alert) (*Alert, error) {
	tx := m.db.Save(&a)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m.GetAlert(a.ID)
}

func (m *Manager) GetAlert(id uint) (*Alert, error) {
	var a Alert

	tx := m.db.Model(&Alert{}).First(&a, id)
	return &a, tx.Error
}

func (m *Manager) GetAlertByUniqueID(id string) (*Alert, error) {
	var a Alert

	m.logger.Infof("looking for any alerts with uid %s", id)
	tx := m.db.Where("unique_identifier = ?", id).Find(&a)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if a.ID > 0 {
		return &a, nil
	}

	return nil, nil
}

func (m *Manager) GetAlertByName(AlertName string) (*Alert, error) {
	var a Alert

	tx := m.db.Where("name = ?", AlertName).First(&a)
	return &a, tx.Error
}

func (m *Manager) GetAlerts() ([]*Alert, error) {
	var alerts []*Alert
	tx := m.db.Find(&alerts)
	return alerts, tx.Error
}

func (m *Manager) DeleteAlert(id uint) error {
	alert, err := m.GetAlert(id)
	if err != nil {
		return err
	}

	if err := alert.SafeToDelete(m); err != nil {
		return err
	}
	tx := m.db.Delete(alert)
	return tx.Error
}

func (m *Manager) GetAlertsForAService(serviceID uint) ([]*Alert, error) {
	var alerts []*Alert
	tx := m.db.Where("service_id = ?", serviceID).Find(&alerts)
	return alerts, tx.Error
}
