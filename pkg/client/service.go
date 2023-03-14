package client

import (
	"fmt"
	"net/http"

	"github.com/klaxn/klaxn-api/pkg/model/service"
)

const (
	serviceSingle = "api/services/%d"
	services      = "api/services"
)

func (c *Client) GetService(id uint) (*service.Service, error) {
	path := fmt.Sprintf(serviceSingle, id)

	var e service.Service
	err := c.getAndParse(path, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *Client) GetServices() ([]*service.Service, error) {
	var e []*service.Service
	err := c.getAndParse(services, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *Client) CreateService(e *service.Service) (*service.Service, error) {
	var ec service.Service
	err := c.postAndParse(services, &e, &ec)
	return &ec, err
}

func (c *Client) DeleteService(id uint) error {
	path := fmt.Sprintf(serviceSingle, id)
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

func (c *Client) UpdateService(e *service.Service) (*service.Service, error) {
	var ec service.Service
	err := c.putAndParse(fmt.Sprintf(serviceSingle, e.ID), &e, &ec)
	return &ec, err
}
