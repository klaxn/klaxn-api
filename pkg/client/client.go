package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/klaxn/klaxn-api/internal/data"
)

type Client struct {
	baseUrl string
	h       *http.Client
}

func New(baseUrl string) (*Client, error) {
	return &Client{
		baseUrl: baseUrl,
		h:       &http.Client{Timeout: time.Second * 10},
	}, nil
}

func (c *Client) do(method, path string, body interface{}) (*http.Response, error) {
	var request *http.Request
	var err error

	url := fmt.Sprintf("%s/%s", c.baseUrl, path)

	if body != nil {
		var requestBytes []byte
		requestBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		request, err = http.NewRequest(method, url, bytes.NewReader(requestBytes))
	} else {
		request, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	response, err := c.h.Do(request)
	if response.StatusCode >= 400 {
		return response, parseErrorResponse(response)
	}

	return response, err
}

func parseErrorResponse(r *http.Response) error {
	var errResponse data.Error
	if err := parseResponse(r, &errResponse); err != nil {
		return fmt.Errorf("got %d status when trying to %s %s", r.StatusCode, r.Request.Method, r.Request.URL.String())
	}
	return fmt.Errorf(errResponse.Message)
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.do(http.MethodDelete, path, nil)
}

func (c *Client) post(path string, body interface{}) (*http.Response, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *Client) put(path string, body interface{}) (*http.Response, error) {
	return c.do(http.MethodPut, path, body)
}

func (c *Client) getAndParse(path string, t interface{}) error {
	response, err := c.get(path)
	if err != nil {
		return err
	}
	return parseResponse(response, t)
}

func (c *Client) postAndParse(path string, body, t interface{}) error {
	response, err := c.post(path, body)
	if err != nil {
		return err
	}
	return parseResponse(response, t)
}

func (c *Client) putAndParse(path string, body, t interface{}) error {
	response, err := c.put(path, body)
	if err != nil {
		return err
	}
	return parseResponse(response, t)
}

func parseResponse(response *http.Response, t interface{}) error {
	defer response.Body.Close()

	responseB, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseB, &t)
}
