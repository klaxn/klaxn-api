package user

import (
	"github.com/klaxn/klaxn-api/internal/data"
	"github.com/klaxn/klaxn-api/pkg/model"
)

type ContactMethod struct {
	model.Item
	Type    string `json:"type"`
	Address string `json:"address"`
	Summary string `json:"summary"`
}

type User struct {
	model.Item
	Email          string          `json:"email"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	UserType       string          `json:"user_type"`
	ContactMethods []ContactMethod `json:"contact_methods,omitempty"`
}

func FromData(d *data.User) *User {
	var cm []ContactMethod
	for _, method := range d.ContactMethods {
		cm = append(cm, ContactMethod{
			Item:    model.Item{ID: method.ID},
			Type:    method.Type,
			Address: method.Address,
			Summary: method.Summary,
		})
	}

	e := &User{
		Item:           model.Item{ID: d.ID},
		Email:          d.Email,
		FirstName:      d.FirstName,
		LastName:       d.LastName,
		UserType:       d.UserType,
		ContactMethods: cm,
	}

	return e
}

func (e *User) ToData() *data.User {
	var cm []data.UserContactMethod
	for _, method := range e.ContactMethods {
		cm = append(cm, data.UserContactMethod{
			Item:    data.Item{ID: method.ID},
			Type:    method.Type,
			Address: method.Address,
			Summary: method.Summary,
		})
	}
	d := &data.User{
		Item:           data.Item{ID: e.ID},
		Email:          e.Email,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		UserType:       e.UserType,
		ContactMethods: cm,
	}

	return d
}

func FromDataSlice(in []*data.User) []*User {
	var out []*User

	for _, service := range in {
		fromData := FromData(service)
		out = append(out, fromData)
	}

	return out
}
