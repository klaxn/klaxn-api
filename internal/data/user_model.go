package data

import (
	"errors"
)

type User struct {
	Item
	Email          string              `json:"email" gorm:"index,unique"`
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	UserType       string              `json:"user_type"`
	ContactMethods []UserContactMethod `json:"contact_methods" gorm:"foreignKey:ID;references:ID"`
}

type UserContactMethod struct {
	Item
	Type    string `json:"type"`
	Address string `json:"address"`
	Summary string `json:"summary"`
}

func (ucm *UserContactMethod) Validate(*Manager) error {
	return nil
}

func (u *User) Validate(*Manager) error {
	if len(u.FirstName) == 0 {
		return errors.New("please set a first name")
	}

	if len(u.LastName) == 0 {
		return errors.New("please set a last name")
	}

	if len(u.Email) == 0 {
		return errors.New("please set an email")
	}

	if len(u.UserType) == 0 {
		return errors.New("please set a user type")
	}

	return nil
}

func (u *User) SafeToDelete(*Manager) error {
	return errors.New("implement me")
}

func (ucm *UserContactMethod) SafeToDelete(*Manager) error {
	return errors.New("implement me")
}
