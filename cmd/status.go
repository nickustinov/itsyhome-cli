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

func showHomeStatus(c *client.Client) error {
	status, err := c.GetStatus()
	if err != nil {
		return err
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	tbl := display.NewTable("", "")
	tbl.AddRow("Rooms", fmt.Sprintf("%d", status.Rooms))
	tbl.AddRow("Devices", fmt.Sprintf("%d", status.Devices))
	tbl.AddRow("Accessories", fmt.Sprintf("%d", status.Accessories))
	tbl.AddRow("Reachable", fmt.Sprintf("%d", status.Reachable))
	tbl.AddRow("Unreachable", fmt.Sprintf("%d", status.Unreachable))
	tbl.AddRow("Scenes", fmt.Sprintf("%d", status.Scenes))
	tbl.AddRow("Groups", fmt.Sprintf("%d", status.Groups))
	fmt.Print(tbl.Render())
	return nil
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
