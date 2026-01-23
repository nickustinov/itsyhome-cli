package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/nickustinov/itsyhome-cli/internal/display"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [room]",
	Short: "Show home status summary, or device states for a room",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())

		if len(args) > 0 {
			return showRoomStatus(c, strings.Join(args, " "))
		}

		return showHomeStatus(c)
	},
}

type statusDevice struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	State string `json:"state"`
	Value string `json:"value,omitempty"`
}

type statusRoom struct {
	Room    string         `json:"room"`
	Devices []statusDevice `json:"devices"`
}

type statusOutput struct {
	Rooms       int          `json:"rooms"`
	Devices     int          `json:"devices"`
	Unreachable int          `json:"unreachable"`
	Details     []statusRoom `json:"details"`
}

func showHomeStatus(c *client.Client) error {
	status, err := c.GetStatus()
	if err != nil {
		return err
	}

	rooms, err := c.ListRooms()
	if err != nil {
		return err
	}

	details := make([]statusRoom, len(rooms))
	maxName, maxType := 0, 0

	for i, room := range rooms {
		infos, err := c.GetInfo(room.Name)
		if err != nil {
			return err
		}
		devices := make([]statusDevice, len(infos))
		for j, info := range infos {
			state := deviceState(info)
			value := ""
			if state == "on" {
				v := formatValue(info)
				if v != "\u2014" {
					value = v
				}
			}
			devices[j] = statusDevice{Name: info.Name, Type: info.Type, State: state, Value: value}
			if len(info.Name) > maxName {
				maxName = len(info.Name)
			}
			if len(info.Type) > maxType {
				maxType = len(info.Type)
			}
		}
		details[i] = statusRoom{Room: room.Name, Devices: devices}
	}

	if jsonOutput {
		out := statusOutput{
			Rooms:       status.Rooms,
			Devices:     status.Devices,
			Unreachable: status.Unreachable,
			Details:     details,
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	header := fmt.Sprintf("Home (%d rooms, %d devices, %d unreachable)",
		status.Rooms, status.Devices, status.Unreachable)

	roomNodes := make([]display.TreeNode, len(rooms))
	for i, room := range rooms {
		deviceNodes := make([]display.TreeNode, len(details[i].Devices))
		for j, dev := range details[i].Devices {
			label := fmt.Sprintf("%-*s  %-*s  %s",
				maxName, dev.Name,
				maxType, dev.Type,
				dev.State)
			if dev.Value != "" {
				label += "    " + dev.Value
			}
			deviceNodes[j] = display.TreeNode{Label: label}
		}
		roomNodes[i] = display.TreeNode{Label: room.Name, Children: deviceNodes}
	}

	tree := &display.Tree{Root: display.TreeNode{Label: header, Children: roomNodes}}
	fmt.Print(tree.Render())
	return nil
}

func deviceState(info client.DeviceInfo) string {
	if !info.Reachable {
		return "unreachable"
	}
	if on, ok := info.State["on"]; ok {
		if b, isBool := on.(bool); isBool && b {
			return "on"
		}
	}
	return "off"
}

func showRoomStatus(c *client.Client, target string) error {
	infos, err := c.GetInfo(target)
	if err != nil {
		return err
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(infos, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	tbl := display.NewTable("Device", "State", "Value")
	for _, info := range infos {
		state := "off"
		value := "\u2014" // em dash

		if on, ok := info.State["on"]; ok {
			if b, isBool := on.(bool); isBool && b {
				state = "on"
			}
		}

		value = formatValue(info)

		tbl.AddRow(info.Name, state, value)
	}
	fmt.Print(tbl.Render())
	return nil
}

func formatValue(info client.DeviceInfo) string {
	parts := []string{}

	if b, ok := info.State["brightness"]; ok {
		parts = append(parts, fmt.Sprintf("%.0f%%", toFloat(b)))
	}
	if t, ok := info.State["temperature"]; ok {
		parts = append(parts, fmt.Sprintf("%.1f\u00b0", toFloat(t)))
	}
	if t, ok := info.State["targetTemperature"]; ok {
		if _, hasCurrent := info.State["temperature"]; !hasCurrent {
			parts = append(parts, fmt.Sprintf("%.1f\u00b0", toFloat(t)))
		}
	}
	if p, ok := info.State["position"]; ok {
		parts = append(parts, fmt.Sprintf("%.0f%%", toFloat(p)))
	}
	if h, ok := info.State["humidity"]; ok {
		parts = append(parts, fmt.Sprintf("%.0f%% RH", toFloat(h)))
	}
	if s, ok := info.State["speed"]; ok {
		parts = append(parts, fmt.Sprintf("speed %.0f%%", toFloat(s)))
	}
	if l, ok := info.State["locked"]; ok {
		if b, isBool := l.(bool); isBool {
			if b {
				parts = append(parts, "locked")
			} else {
				parts = append(parts, "unlocked")
			}
		}
	}

	if len(parts) == 0 {
		return "\u2014"
	}
	return strings.Join(parts, ", ")
}

func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case json.Number:
		f, _ := n.Float64()
		return f
	}
	return 0
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
