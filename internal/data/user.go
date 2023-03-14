package data

func (m *Manager) CreateUser(u *User) (*User, error) {
	tx := m.db.Create(&u)
	return u, tx.Error
}

func (m *Manager) UpdateUser(u *User) (*User, error) {
	tx := m.db.Save(&u)
	return u, tx.Error
}

func (m *Manager) GetUser(id uint) (*User, error) {
	var u User

	tx := m.db.Model(&User{}).First(&u, id)
	return &u, tx.Error
}

func (m *Manager) GetUserByEmail(email string) (*User, error) {
	var u User

	tx := m.db.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &u, nil
}

func (m *Manager) GetUsers() ([]*User, error) {
	var users []*User
	tx := m.db.Find(&users)
	return users, tx.Error
}

func (m *Manager) DeleteUser(id uint) error {
	user, err := m.GetUser(id)
	if err != nil {
		return err
	}

	if err := user.SafeToDelete(m); err != nil {
		return err
	}
	tx := m.db.Delete(user)
	return tx.Error
}
