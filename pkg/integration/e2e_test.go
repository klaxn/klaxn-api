package integration

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/klaxn/klaxn-api/pkg/model/escalation"
	"github.com/klaxn/klaxn-api/pkg/model/service"
	"github.com/klaxn/klaxn-api/pkg/model/team"
)

func TestEscalations(t *testing.T) {
	c := setup(t)

	t.Run("create query and delete", func(tt *testing.T) {
		escalationName := "e2e test escalation 1"
		escalationDescription := "This is a description"
		newEscalation, err := c.CreateEscalation(&escalation.Escalation{
			Name:        escalationName,
			Description: escalationDescription,
			Layers:      []escalation.Layer{{Tier: 1, ResponderType: "user", ResponderReference: "h"}},
		})
		if err != nil {
			tt.Fatal(err)
		}

		assert.Equal(tt, escalationName, newEscalation.Name)
		assert.Equal(tt, escalationDescription, newEscalation.Description)

		getEscalation, err := c.GetEscalation(newEscalation.ID)
		if err != nil {
			tt.Fatal(err)
		}
		assert.Equal(tt, escalationName, getEscalation.Name)
		assert.Equal(tt, escalationDescription, getEscalation.Description)
		assert.Equal(tt, 1, len(getEscalation.Layers))

		escalations, err := c.GetEscalations()
		if err != nil {
			tt.Fatal(err)
		}

		assert.Equal(tt, 1, len(escalations))

		err = c.DeleteEscalation(escalations[0].ID)
		if err != nil {
			tt.Fatal(err)
		}
	})
	t.Run("create and update", func(tt *testing.T) {
		escalationName := "e2e test escalation 2"
		escalationDescription1 := "This is a description 1st version"
		escalationDescription2 := "This is a description 2nd version"
		newEscalation, err := c.CreateEscalation(&escalation.Escalation{
			Name:        escalationName,
			Description: escalationDescription1,
			Layers:      []escalation.Layer{{Tier: 1, ResponderType: "user", ResponderReference: "h"}},
		})
		if err != nil {
			tt.Fatal(err)
		}

		get1, err := c.GetEscalation(newEscalation.ID)
		if err != nil {
			tt.Fatal(err)
		}
		assert.Equal(tt, escalationDescription1, get1.Description)

		get1.Description = escalationDescription2
		update1, err := c.UpdateEscalation(get1)
		if err != nil {
			tt.Fatal(err)
		}
		assert.Equal(tt, escalationDescription2, update1.Description)

	})
}

func TestServices(t *testing.T) {
	c := setup(t)

	t.Run("create", func(tt *testing.T) {
		escalationName := "e2e test escalation 1"
		escalationDescription := "This is an escalation description"
		e, err := c.CreateEscalation(&escalation.Escalation{
			Name:        escalationName,
			Description: escalationDescription,
			Layers:      []escalation.Layer{{Tier: 1, ResponderType: "user", ResponderReference: "h"}},
		})
		if err != nil {
			tt.Fatal(err)
		}

		teamName := "e2e test team 1"
		teamDescription := "This is a team description"
		te, err := c.CreateTeam(&team.Team{
			Name:        teamName,
			Description: teamDescription,
		})
		if err != nil {
			tt.Fatal(err)
		}

		serviceName := "e2e test service 1"
		serviceDescription := "This is a service description"
		createService, err := c.CreateService(&service.Service{
			Name:         serviceName,
			TeamID:       te.ID,
			Description:  serviceDescription,
			EscalationID: e.ID,
		})
		if err != nil {
			tt.Fatal(err)
		}

		assert.Equal(tt, serviceName, createService.Name)
		assert.Equal(tt, serviceDescription, createService.Description)
		assert.Equal(tt, te.ID, createService.TeamID)
		assert.Equal(tt, e.ID, createService.EscalationID)
	})
}
