package client

import (
	"fmt"
	"net/http"

	"github.com/klaxn/klaxn-api/pkg/model/team"
)

const (
	teamSingle = "api/teams/%d"
	teams      = "api/teams"
)

func (c *Client) GetTeam(id uint) (*team.Team, error) {
	path := fmt.Sprintf(teamSingle, id)

	var e team.Team
	err := c.getAndParse(path, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *Client) GetTeams() ([]*team.Team, error) {
	var e []*team.Team
	err := c.getAndParse(teams, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *Client) CreateTeam(t *team.Team) (*team.Team, error) {
	var ec team.Team
	err := c.postAndParse(teams, &t, &ec)
	return &ec, err
}

func (c *Client) DeleteTeam(id uint) error {
	path := fmt.Sprintf(teamSingle, id)
	response, err := c.delete(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent {
		return nil
	}

	return parseErrorResponse(response)
}

func (c *Client) UpdateTeam(t *team.Team) (*team.Team, error) {
	var ec team.Team
	err := c.putAndParse(fmt.Sprintf(teamSingle, t.ID), &t, &ec)
	return &ec, err
}
