package data

import (
	"errors"
	"fmt"
)

type Team struct {
	Item
	Name        string `json:"name" gorm:"index,unique"`
	Description string `json:"description"`
}

func (t *Team) Validate(*Manager) error {
	fmt.Println("WARNING Implement model.Team.Validate()")
	return nil
}

func (t *Team) SafeToDelete(*Manager) error {
	return errors.New("implement me")
}
