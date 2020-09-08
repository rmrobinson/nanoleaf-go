package nanoleaf

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	// ErrUnauthorized is returned if the specified API key is invalid
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden is returned if the specified API key is invalid
	ErrForbidden = errors.New("forbidden")
	// ErrBadRequest is returned if the request had malformed data
	ErrBadRequest = errors.New("bad request")
	// ErrNotFound is returned if the resource was not found
	ErrNotFound = errors.New("not found")
	// ErrUnknown is returned if the specific failure reason isn't known
	ErrUnknown = errors.New("unknown")
)

// Client represents a handle to the Nanoleaf API
type Client struct {
	httpClient *http.Client
	apiKey     string
	hostname   string
	port       int
}

// NewClient creates a new Nanoleaf API client
func NewClient(httpClient *http.Client, hostname string, port int, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		hostname:   hostname,
		port:       port,
		apiKey:     apiKey,
	}
}

func (c *Client) getURLBase() string {
	return "http://" + c.hostname + ":" + strconv.Itoa(c.port) + "/api/v1/" + c.apiKey + "/"
}

// CreateAPIKey is used to create a new API key (after pressing the 'pair' button on the panel)
func (c *Client) CreateAPIKey(ctx context.Context) (string, error) {
	r, err := http.NewRequest(http.MethodPost, "http://"+c.hostname+":"+strconv.Itoa(c.port)+"/api/v1/new", nil)
	if err != nil {
		return "", err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var nanoleafResp struct {
		AuthToken string `json:"auth_token"`
	}

	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&nanoleafResp)
		if err != nil {
			return "", err
		}

		return nanoleafResp.AuthToken, nil
	} else if resp.StatusCode == 401 {
		return "", ErrUnauthorized
	} else if resp.StatusCode == 403 {
		return "", ErrForbidden
	} else if resp.StatusCode == 422 {
		return "", ErrBadRequest
	}

	return "", ErrUnknown
}

// DeleteAPIKey is used to remove an existing API key
func (c *Client) DeleteAPIKey(ctx context.Context, keyToDelete string) error {
	r, err := http.NewRequest(http.MethodDelete, c.getURLBase(), nil)
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	} else if resp.StatusCode == 401 {
		return ErrUnauthorized
	}

	return ErrUnknown
}

func (c *Client) get(ctx context.Context, path string, respType interface{}) error {
	r, err := http.NewRequest(http.MethodGet, c.getURLBase()+path, nil)
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return json.NewDecoder(resp.Body).Decode(respType)
	} else if resp.StatusCode == 400 {
		return ErrBadRequest
	} else if resp.StatusCode == 401 {
		return ErrUnauthorized
	} else if resp.StatusCode == 404 {
		return ErrNotFound
	}

	return ErrUnknown
}

func (c *Client) put(ctx context.Context, path string, reqType interface{}) error {
	req, err := json.Marshal(reqType)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPut, c.getURLBase()+path, bytes.NewBuffer(req))
	if err != nil {
		return err
	}

	r = r.WithContext(ctx)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	} else if resp.StatusCode == 400 {
		return ErrBadRequest
	} else if resp.StatusCode == 401 {
		return ErrUnauthorized
	} else if resp.StatusCode == 404 {
		return ErrNotFound
	} else if resp.StatusCode == 422 {
		return ErrBadRequest
	}

	return ErrUnknown
}
