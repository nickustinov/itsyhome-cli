package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nickustinov/itsyhome-cli/internal/client"
)

// setupTestEnv sets up a test server and config pointing to it.
// The config.json will have host and port matching the test server.
func setupTestEnv(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "itsyhome")
	os.MkdirAll(dir, 0755)

	// The BaseURL format is http://host:port, parse the listener address
	addr := srv.Listener.Addr().String()
	// addr is like "127.0.0.1:PORT" or "[::1]:PORT"
	// We write the config so that BaseURL() produces the correct URL
	// config.BaseURL() returns fmt.Sprintf("http://%s:%d", c.Host, c.Port)
	// Since httptest URL is http://127.0.0.1:PORT, we need host=127.0.0.1, port=PORT
	var host string
	var port int
	// Parse "host:port" from addr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			host = addr[:i]
			fmt.Sscanf(addr[i+1:], "%d", &port)
			break
		}
	}

	cfgData, _ := json.Marshal(map[string]interface{}{
		"host": host,
		"port": port,
	})
	os.WriteFile(filepath.Join(dir, "config.json"), cfgData, 0644)
}

func executeCmd(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

// --- status command tests ---

func TestStatusCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			json.NewEncoder(w).Encode(map[string]int{
				"rooms": 2, "devices": 3, "accessories": 5,
				"reachable": 2, "unreachable": 1, "scenes": 4, "groups": 2,
			})
		case "/list/rooms":
			json.NewEncoder(w).Encode([]map[string]string{
				{"name": "Office"},
				{"name": "Bedroom"},
			})
		case "/info/Office":
			json.NewEncoder(w).Encode([]client.DeviceInfo{
				{Name: "Desk Lamp", Type: "light", Reachable: true, State: map[string]interface{}{"on": true, "brightness": float64(80)}},
				{Name: "AC Unit", Type: "thermostat", Reachable: true, State: map[string]interface{}{"on": true, "temperature": float64(22.5)}},
			})
		case "/info/Bedroom":
			json.NewEncoder(w).Encode([]client.DeviceInfo{
				{Name: "Floor Lamp", Type: "light", Reachable: false, State: map[string]interface{}{}},
			})
		}
	})

	jsonOutput = false
	_, err := executeCmd("status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdEmptyRoom(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			json.NewEncoder(w).Encode(map[string]int{
				"rooms": 1, "devices": 0, "accessories": 0,
				"reachable": 0, "unreachable": 0, "scenes": 0, "groups": 0,
			})
		case "/list/rooms":
			json.NewEncoder(w).Encode([]map[string]string{{"name": "Empty Room"}})
		case "/info/Empty Room":
			json.NewEncoder(w).Encode([]client.DeviceInfo{})
		}
	})

	jsonOutput = false
	_, err := executeCmd("status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdAllUnreachable(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/status":
			json.NewEncoder(w).Encode(map[string]int{
				"rooms": 1, "devices": 2, "accessories": 0,
				"reachable": 0, "unreachable": 2, "scenes": 0, "groups": 0,
			})
		case "/list/rooms":
			json.NewEncoder(w).Encode([]map[string]string{{"name": "Office"}})
		case "/info/Office":
			json.NewEncoder(w).Encode([]client.DeviceInfo{
				{Name: "Lamp", Type: "light", Reachable: false, State: map[string]interface{}{}},
				{Name: "Fan", Type: "fan", Reachable: false, State: map[string]interface{}{}},
			})
		}
	})

	jsonOutput = false
	_, err := executeCmd("status")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]int{
			"rooms": 3, "devices": 10, "accessories": 5,
			"reachable": 8, "unreachable": 2, "scenes": 4, "groups": 2,
		})
	})

	jsonOutput = false
	_, err := executeCmd("status", "--json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdWithRoom(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]client.DeviceInfo{
			{Name: "Lamp", Type: "light", Reachable: true, State: map[string]interface{}{"on": true, "brightness": float64(80)}},
		})
	})

	jsonOutput = false
	_, err := executeCmd("status", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdWithRoomJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]client.DeviceInfo{
			{Name: "Lamp", Type: "light", Reachable: true, State: map[string]interface{}{"on": true}},
		})
	})

	jsonOutput = false
	_, err := executeCmd("status", "--json", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("status")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStatusCmdRoomError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("status", "Office")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- formatValue tests ---

func TestFormatValueBrightness(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"brightness": float64(75)}}
	result := formatValue(info)
	if result != "75%" {
		t.Errorf("expected 75%%, got %s", result)
	}
}

func TestFormatValueTemperature(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"temperature": float64(22.5)}}
	result := formatValue(info)
	if result != "22.5\u00b0" {
		t.Errorf("expected 22.5°, got %s", result)
	}
}

func TestFormatValueTargetTemperature(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"targetTemperature": float64(21.0)}}
	result := formatValue(info)
	if result != "21.0\u00b0" {
		t.Errorf("expected 21.0°, got %s", result)
	}
}

func TestFormatValueTargetTemperatureWithCurrent(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{
		"temperature":       float64(22.5),
		"targetTemperature": float64(21.0),
	}}
	result := formatValue(info)
	// Should only show current temperature, not target
	if result != "22.5\u00b0" {
		t.Errorf("expected 22.5°, got %s", result)
	}
}

func TestFormatValuePosition(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"position": float64(50)}}
	result := formatValue(info)
	if result != "50%" {
		t.Errorf("expected 50%%, got %s", result)
	}
}

func TestFormatValueHumidity(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"humidity": float64(60)}}
	result := formatValue(info)
	if result != "60% RH" {
		t.Errorf("expected 60%% RH, got %s", result)
	}
}

func TestFormatValueSpeed(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"speed": float64(80)}}
	result := formatValue(info)
	if result != "speed 80%" {
		t.Errorf("expected speed 80%%, got %s", result)
	}
}

func TestFormatValueLocked(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"locked": true}}
	result := formatValue(info)
	if result != "locked" {
		t.Errorf("expected locked, got %s", result)
	}
}

func TestFormatValueUnlocked(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"locked": false}}
	result := formatValue(info)
	if result != "unlocked" {
		t.Errorf("expected unlocked, got %s", result)
	}
}

func TestFormatValueLockedNonBool(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{"locked": "yes"}}
	result := formatValue(info)
	if result != "\u2014" {
		t.Errorf("expected em dash, got %s", result)
	}
}

func TestFormatValueEmpty(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{}}
	result := formatValue(info)
	if result != "\u2014" {
		t.Errorf("expected em dash, got %s", result)
	}
}

func TestFormatValueMultiple(t *testing.T) {
	info := client.DeviceInfo{State: map[string]interface{}{
		"brightness": float64(75),
		"position":   float64(50),
	}}
	result := formatValue(info)
	// Both should appear
	if result != "75%, 50%" {
		t.Errorf("expected '75%%, 50%%', got %s", result)
	}
}

// --- toFloat tests ---

func TestToFloatFloat64(t *testing.T) {
	if toFloat(float64(42.5)) != 42.5 {
		t.Error("expected 42.5")
	}
}

func TestToFloatInt(t *testing.T) {
	if toFloat(int(10)) != 10.0 {
		t.Error("expected 10.0")
	}
}

func TestToFloatJSONNumber(t *testing.T) {
	n := json.Number("3.14")
	if toFloat(n) != 3.14 {
		t.Errorf("expected 3.14, got %f", toFloat(n))
	}
}

func TestToFloatUnknown(t *testing.T) {
	if toFloat("hello") != 0 {
		t.Error("expected 0 for unknown type")
	}
}

// --- control command tests ---

func TestToggleCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("toggle", "Office", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestToggleCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("toggle", "--json", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestToggleCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "not found"})
	})

	jsonOutput = false
	_, err := executeCmd("toggle", "Unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOnCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("on", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOffCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("off", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBrightnessCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("brightness", "80", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSceneCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	jsonOutput = false
	_, err := executeCmd("scene", "Goodnight")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- info command tests ---

func TestInfoCmdSingle(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(client.DeviceInfo{
			Name: "Lamp", Type: "light", Room: "Office", Reachable: true,
			State: map[string]interface{}{"on": true, "brightness": float64(80)},
		})
	})

	jsonOutput = false
	_, err := executeCmd("info", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfoCmdSingleUnreachable(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(client.DeviceInfo{
			Name: "Lamp", Type: "light", Reachable: false,
		})
	})

	jsonOutput = false
	_, err := executeCmd("info", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfoCmdMultiple(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]client.DeviceInfo{
			{Name: "Lamp 1", Type: "light", Reachable: true, State: map[string]interface{}{"on": true}},
			{Name: "Lamp 2", Type: "light", Reachable: false, State: map[string]interface{}{"on": false}},
		})
	})

	jsonOutput = false
	_, err := executeCmd("info", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfoCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(client.DeviceInfo{
			Name: "Lamp", Type: "light", Reachable: true,
		})
	})

	jsonOutput = false
	_, err := executeCmd("info", "--json", "Lamp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfoCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("info", "Unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- list command tests ---

func TestListRoomsCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]string{{"name": "Office"}, {"name": "Bedroom"}})
	})

	jsonOutput = false
	_, err := executeCmd("list", "rooms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRoomsCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]string{{"name": "Office"}})
	})

	jsonOutput = false
	_, err := executeCmd("list", "--json", "rooms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRoomsCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("list", "rooms")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListDevicesCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": "Lamp", "type": "light", "room": "Office", "reachable": true},
			{"name": "Fan", "type": "fan", "room": "Bedroom", "reachable": false},
		})
	})

	jsonOutput = false
	_, err := executeCmd("list", "devices")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListDevicesCmdWithRoom(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": "Lamp", "type": "light", "room": "Office", "reachable": true},
		})
	})

	jsonOutput = false
	_, err := executeCmd("list", "devices", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListDevicesCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": "Lamp", "type": "light", "reachable": true},
		})
	})

	jsonOutput = false
	_, err := executeCmd("list", "--json", "devices")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListDevicesCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("list", "devices")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListScenesCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]string{{"name": "Goodnight"}})
	})

	jsonOutput = false
	_, err := executeCmd("list", "scenes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListScenesCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]string{{"name": "Goodnight"}})
	})

	jsonOutput = false
	_, err := executeCmd("list", "--json", "scenes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListScenesCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("list", "scenes")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListGroupsCmd(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": "All Lights", "icon": "💡", "devices": 5},
		})
	})

	jsonOutput = false
	_, err := executeCmd("list", "groups")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListGroupsCmdJSON(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": "All Lights", "icon": "💡", "devices": 5},
		})
	})

	jsonOutput = false
	_, err := executeCmd("list", "--json", "groups")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListGroupsCmdError(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	jsonOutput = false
	_, err := executeCmd("list", "groups")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- config command tests ---

func TestConfigCmd(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	jsonOutput = false
	_, err := executeCmd("config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigSetCmd(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	jsonOutput = false
	_, err := executeCmd("config", "set", "--host", "10.0.0.1", "--port", "9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it was saved
	_, err = executeCmd("config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigSetCmdSaveError(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Make the config directory a file to cause a save error
	configDir := filepath.Join(tmp, ".config", "itsyhome")
	os.MkdirAll(filepath.Dir(configDir), 0755)
	os.WriteFile(configDir, []byte("not a dir"), 0644)

	jsonOutput = false
	// This should print an error but not return one (Run not RunE)
	_, err := executeCmd("config", "set", "--host", "10.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- root command tests ---

func TestRootCmdHelp(t *testing.T) {
	jsonOutput = false
	_, err := executeCmd("--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCmdUnknown(t *testing.T) {
	jsonOutput = false
	_, err := executeCmd("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestExecuteSuccess(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	Execute()
}

func TestExecuteError(t *testing.T) {
	var exitCode int
	original := osExit
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = original }()

	rootCmd.SetArgs([]string{"nonexistent"})
	Execute()

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

// --- showRoomStatus with device states ---

func TestShowRoomStatusDeviceOff(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]client.DeviceInfo{
			{Name: "Lamp", Type: "light", Reachable: true, State: map[string]interface{}{"on": false}},
		})
	})

	jsonOutput = false
	_, err := executeCmd("status", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShowRoomStatusNoState(t *testing.T) {
	setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]client.DeviceInfo{
			{Name: "Sensor", Type: "sensor", Reachable: true, State: map[string]interface{}{}},
		})
	})

	jsonOutput = false
	_, err := executeCmd("status", "Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
