package data

func (m *Manager) GetTeam(id uint) (*Team, error) {
	var t Team

	tx := m.db.Model(&Team{}).First(&t, id)
	return &t, tx.Error
}

func (m *Manager) GetTeamByName(teamName string) (*Team, error) {
	var t Team

	tx := m.db.Where("name = ?", teamName).First(&t)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &t, nil
}

func (m *Manager) GetTeams() ([]*Team, error) {
	var teams []*Team
	tx := m.db.Find(&teams)
	return teams, tx.Error
}

func (m *Manager) UpdateTeam(t *Team) (*Team, error) {
	tx := m.db.Save(&t)
	return t, tx.Error
}

func (m *Manager) DeleteTeam(id uint) error {
	team, err := m.GetTeam(id)
	if err != nil {
		return err
	}

	if err := team.SafeToDelete(m); err != nil {
		return err
	}
	tx := m.db.Delete(team)
	return tx.Error
}
