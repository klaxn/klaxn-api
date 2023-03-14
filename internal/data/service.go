package data

func (m *Manager) UpdateService(s *Service) (*Service, error) {
	tx := m.db.Save(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return m.GetService(s.ID)
}

func (m *Manager) GetService(id uint) (*Service, error) {
	var s Service

	tx := m.db.Model(&Service{}).First(&s, id)
	return &s, tx.Error
}

func (m *Manager) GetServiceByName(serviceName string) (*Service, error) {
	var s Service

	tx := m.db.Where("name = ?", serviceName).First(&s)
	return &s, tx.Error
}

func (m *Manager) GetServices() ([]*Service, error) {
	var services []*Service
	tx := m.db.Find(&services)
	return services, tx.Error
}

func (m *Manager) DeleteService(id uint) error {
	service, err := m.GetService(id)
	if err != nil {
		return err
	}

	if err := service.SafeToDelete(m); err != nil {
		return err
	}
	tx := m.db.Delete(service)
	return tx.Error
}

func (m *Manager) GetServicesOnEscalation(escalationID uint) ([]*Service, error) {
	var services []*Service
	tx := m.db.Where("escalation_id = ?", escalationID).Find(&services)
	return services, tx.Error
}
