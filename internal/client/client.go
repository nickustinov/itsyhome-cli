package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nickustinov/itsyhome-cli/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ActionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type StatusResponse struct {
	Rooms       int `json:"rooms"`
	Devices     int `json:"devices"`
	Accessories int `json:"accessories"`
	Reachable   int `json:"reachable"`
	Unreachable int `json:"unreachable"`
	Scenes      int `json:"scenes"`
	Groups      int `json:"groups"`
}

type Room struct {
	Name string `json:"name"`
}

type Device struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Room      string `json:"room,omitempty"`
	Reachable bool   `json:"reachable"`
}

type Scene struct {
	Name string `json:"name"`
}

type Group struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Devices int    `json:"devices"`
	Room    string `json:"room,omitempty"`
}

type DeviceInfo struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Room      string                 `json:"room,omitempty"`
	Reachable bool                   `json:"reachable"`
	State     map[string]interface{} `json:"state,omitempty"`
}

func New(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.BaseURL(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) DoAction(path string) (*ActionResponse, error) {
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var resp ActionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if resp.Status == "error" {
		return nil, fmt.Errorf("%s", resp.Message)
	}

	return &resp, nil
}

func (c *Client) GetStatus() (*StatusResponse, error) {
	body, err := c.get("/status")
	if err != nil {
		return nil, err
	}

	var resp StatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &resp, nil
}

func (c *Client) ListRooms() ([]Room, error) {
	body, err := c.get("/list/rooms")
	if err != nil {
		return nil, err
	}

	var rooms []Room
	if err := json.Unmarshal(body, &rooms); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return rooms, nil
}

func (c *Client) ListDevices(room string) ([]Device, error) {
	path := "/list/devices"
	if room != "" {
		path += "/" + url.PathEscape(room)
	}

	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var devices []Device
	if err := json.Unmarshal(body, &devices); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return devices, nil
}

func (c *Client) ListScenes() ([]Scene, error) {
	body, err := c.get("/list/scenes")
	if err != nil {
		return nil, err
	}

	var scenes []Scene
	if err := json.Unmarshal(body, &scenes); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return scenes, nil
}

func (c *Client) ListGroups() ([]Group, error) {
	body, err := c.get("/list/groups")
	if err != nil {
		return nil, err
	}

	var groups []Group
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return groups, nil
}

func (c *Client) GetInfo(target string) ([]DeviceInfo, error) {
	path := "/info/" + url.PathEscape(target)

	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	// Check for error response first
	var errResp ActionResponse
	if json.Unmarshal(body, &errResp) == nil && errResp.Status == "error" {
		return nil, fmt.Errorf("%s", errResp.Message)
	}

	// Try array
	var infos []DeviceInfo
	if err := json.Unmarshal(body, &infos); err == nil {
		return infos, nil
	}

	// Try single object
	var info DeviceInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return []DeviceInfo{info}, nil
}

func (c *Client) get(path string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w\nIs the Itsyhome app running with the server enabled?\nNote: webhook/CLI access requires an Itsyhome Pro subscription.", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("Itsyhome Pro required for webhook/CLI access")
	}

	if resp.StatusCode >= 400 {
		var errResp ActionResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
			return nil, fmt.Errorf("%s", errResp.Message)
		}
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	return body, nil
}
