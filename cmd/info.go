package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/nickustinov/itsyhome-cli/internal/display"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <device|room|group>",
	Short: "Show detailed info about a device, room, or group",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := strings.Join(args, " ")
		c := client.New(config.Load())
		infos, err := c.GetInfo(target)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(infos, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(infos) == 1 {
			printSingleInfo(infos[0])
		} else {
			printMultiInfo(infos)
		}
		return nil
	},
}

func printSingleInfo(info client.DeviceInfo) {
	tbl := display.NewTable("Property", "Value")
	tbl.AddRow("Name", info.Name)
	tbl.AddRow("Type", info.Type)
	if info.Room != "" {
		tbl.AddRow("Room", info.Room)
	}
	if info.Reachable {
		tbl.AddRow("Status", "reachable")
	} else {
		tbl.AddRow("Status", "unreachable")
	}

	if len(info.State) > 0 {
		keys := make([]string, 0, len(info.State))
		for k := range info.State {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			tbl.AddRow(k, fmt.Sprintf("%v", info.State[k]))
		}
	}
	fmt.Print(tbl.Render())
}

func printMultiInfo(infos []client.DeviceInfo) {
	tbl := display.NewTable("Device", "Type", "State", "Value")
	for _, info := range infos {
		state := "off"
		if on, ok := info.State["on"]; ok {
			if b, isBool := on.(bool); isBool && b {
				state = "on"
			}
		}
		if !info.Reachable {
			state = "unreachable"
		}
		tbl.AddRow(info.Name, info.Type, state, formatValue(info))
	}
	fmt.Print(tbl.Render())
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
