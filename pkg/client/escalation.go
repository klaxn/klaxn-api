package client

import (
	"fmt"
	"net/http"

	"github.com/klaxn/klaxn-api/pkg/model/escalation"
)

const (
	escalationSingle = "api/escalations/%d"
	escalations      = "api/escalations"
)

func (c *Client) GetEscalation(id uint) (*escalation.Escalation, error) {
	path := fmt.Sprintf(escalationSingle, id)

	var e escalation.Escalation
	err := c.getAndParse(path, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *Client) GetEscalations() ([]*escalation.Escalation, error) {
	var e []*escalation.Escalation
	err := c.getAndParse(escalations, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *Client) CreateEscalation(e *escalation.Escalation) (*escalation.Escalation, error) {
	var ec escalation.Escalation
	err := c.postAndParse(escalations, &e, &ec)
	return &ec, err
}

func (c *Client) DeleteEscalation(id uint) error {
	path := fmt.Sprintf(escalationSingle, id)
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

func (c *Client) UpdateEscalation(e *escalation.Escalation) (*escalation.Escalation, error) {
	var ec escalation.Escalation
	err := c.putAndParse(fmt.Sprintf(escalationSingle, e.ID), &e, &ec)
	return &ec, err
}
