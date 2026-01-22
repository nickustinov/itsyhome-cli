package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickustinov/itsyhome-cli/internal/config"
)

func testServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	srv := httptest.NewServer(handler)
	c := &Client{
		baseURL:    srv.URL,
		httpClient: srv.Client(),
	}
	return srv, c
}

func TestDoAction(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/toggle/Office/Lamp" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(ActionResponse{Status: "success"})
	})
	defer srv.Close()

	resp, err := c.DoAction("/toggle/Office/Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}

func TestDoActionError(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ActionResponse{Status: "error", Message: "device not found"})
	})
	defer srv.Close()

	_, err := c.DoAction("/toggle/Unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "device not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestGetStatus(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(StatusResponse{
			Rooms: 3, Devices: 10, Accessories: 5,
			Reachable: 8, Unreachable: 2, Scenes: 4, Groups: 2,
		})
	})
	defer srv.Close()

	status, err := c.GetStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Rooms != 3 {
		t.Errorf("expected 3 rooms, got %d", status.Rooms)
	}
	if status.Devices != 10 {
		t.Errorf("expected 10 devices, got %d", status.Devices)
	}
}

func TestListRooms(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/list/rooms" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Room{{Name: "Office"}, {Name: "Bedroom"}})
	})
	defer srv.Close()

	rooms, err := c.ListRooms()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rooms) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(rooms))
	}
	if rooms[0].Name != "Office" {
		t.Errorf("expected Office, got %s", rooms[0].Name)
	}
}

func TestListDevices(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/list/devices" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Device{
			{Name: "Lamp", Type: "light", Room: "Office", Reachable: true},
		})
	})
	defer srv.Close()

	devices, err := c.ListDevices("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	if devices[0].Name != "Lamp" {
		t.Errorf("expected Lamp, got %s", devices[0].Name)
	}
}

func TestListDevicesWithRoom(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/list/devices/Office" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]Device{})
	})
	defer srv.Close()

	_, err := c.ListDevices("Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListScenes(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Scene{{Name: "Goodnight"}, {Name: "Morning"}})
	})
	defer srv.Close()

	scenes, err := c.ListScenes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scenes) != 2 {
		t.Fatalf("expected 2 scenes, got %d", len(scenes))
	}
}

func TestListGroups(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Group{
			{Name: "All Lights", Icon: "💡", Devices: 5},
		})
	})
	defer srv.Close()

	groups, err := c.ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Devices != 5 {
		t.Errorf("expected 5 devices, got %d", groups[0].Devices)
	}
}

func TestGetInfo(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/info/Office Lamp" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(DeviceInfo{
			Name: "Office Lamp", Type: "light", Room: "Office", Reachable: true,
			State: map[string]interface{}{"on": true, "brightness": float64(80)},
		})
	})
	defer srv.Close()

	infos, err := c.GetInfo("Office Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(infos) != 1 {
		t.Fatalf("expected 1 info, got %d", len(infos))
	}
	if infos[0].Name != "Office Lamp" {
		t.Errorf("expected Office Lamp, got %s", infos[0].Name)
	}
	if infos[0].State["on"] != true {
		t.Error("expected on=true")
	}
}

func TestGetInfoArray(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]DeviceInfo{
			{Name: "Lamp 1", Type: "light", Reachable: true},
			{Name: "Lamp 2", Type: "light", Reachable: false},
		})
	})
	defer srv.Close()

	infos, err := c.GetInfo("Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(infos) != 2 {
		t.Fatalf("expected 2 infos, got %d", len(infos))
	}
}

func TestGetInfoError(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ActionResponse{Status: "error", Message: "not found"})
	})
	defer srv.Close()

	_, err := c.GetInfo("Unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHTTP403(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	})
	defer srv.Close()

	_, err := c.GetStatus()
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if err.Error() != "Itsyhome Pro required for webhook/CLI access" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestHTTP500(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"status":"error","message":"internal error"}`))
	})
	defer srv.Close()

	_, err := c.GetStatus()
	if err == nil {
		t.Fatal("expected error for 500")
	}
	if err.Error() != "internal error" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestHTTP500NoBody(t *testing.T) {
	srv, c := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
	defer srv.Close()

	_, err := c.GetStatus()
	if err == nil {
		t.Fatal("expected error for 500")
	}
}

func TestNewClient(t *testing.T) {
	cfg := config.Config{Host: "example.com", Port: 1234}
	c := New(cfg)
	if c.baseURL != "http://example.com:1234" {
		t.Errorf("unexpected baseURL: %s", c.baseURL)
	}
}

func TestConnectionRefused(t *testing.T) {
	cfg := config.Config{Host: "localhost", Port: 1}
	c := New(cfg)
	_, err := c.GetStatus()
	if err == nil {
		t.Fatal("expected connection error")
	}
}
