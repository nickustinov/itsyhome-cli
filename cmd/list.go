package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/nickustinov/itsyhome-cli/internal/display"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List rooms, devices, scenes, or groups",
}

var listRoomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "List all rooms",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())
		rooms, err := c.ListRooms()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(rooms, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		tbl := display.NewTable("Room")
		for _, r := range rooms {
			tbl.AddRow(r.Name)
		}
		fmt.Print(tbl.Render())
		return nil
	},
}

var listDevicesCmd = &cobra.Command{
	Use:   "devices [room]",
	Short: "List devices, optionally filtered by room",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())
		room := ""
		if len(args) > 0 {
			room = args[0]
		}

		devices, err := c.ListDevices(room)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(devices, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		tbl := display.NewTable("Device", "Type", "Room", "Status")
		for _, d := range devices {
			status := "ok"
			if !d.Reachable {
				status = "unreachable"
			}
			tbl.AddRow(d.Name, d.Type, d.Room, status)
		}
		fmt.Print(tbl.Render())
		return nil
	},
}

var listScenesCmd = &cobra.Command{
	Use:   "scenes",
	Short: "List all scenes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())
		scenes, err := c.ListScenes()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(scenes, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		tbl := display.NewTable("Scene")
		for _, s := range scenes {
			tbl.AddRow(s.Name)
		}
		fmt.Print(tbl.Render())
		return nil
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())
		groups, err := c.ListGroups()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(groups, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		tbl := display.NewTable("Group", "Icon", "Devices")
		for _, g := range groups {
			tbl.AddRow(g.Name, g.Icon, fmt.Sprintf("%d", g.Devices))
		}
		fmt.Print(tbl.Render())
		return nil
	},
}

func init() {
	listCmd.AddCommand(listRoomsCmd)
	listCmd.AddCommand(listDevicesCmd)
	listCmd.AddCommand(listScenesCmd)
	listCmd.AddCommand(listGroupsCmd)
	rootCmd.AddCommand(listCmd)
}
