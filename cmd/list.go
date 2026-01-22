package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
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

		for _, r := range rooms {
			fmt.Println(r.Name)
		}
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

		for _, d := range devices {
			status := "+"
			if !d.Reachable {
				status = "-"
			}
			if d.Room != "" {
				fmt.Printf("%s %s (%s) [%s]\n", status, d.Name, d.Type, d.Room)
			} else {
				fmt.Printf("%s %s (%s)\n", status, d.Name, d.Type)
			}
		}
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

		for _, s := range scenes {
			fmt.Println(s.Name)
		}
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

		for _, g := range groups {
			fmt.Printf("%s %s (%d devices)\n", g.Icon, g.Name, g.Devices)
		}
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
