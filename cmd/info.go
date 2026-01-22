package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
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

		for i, info := range infos {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("Name: %s\n", info.Name)
			fmt.Printf("Type: %s\n", info.Type)
			if info.Room != "" {
				fmt.Printf("Room: %s\n", info.Room)
			}
			if info.Reachable {
				fmt.Println("Status: reachable")
			} else {
				fmt.Println("Status: unreachable")
			}

			if len(info.State) > 0 {
				fmt.Println("State:")
				keys := make([]string, 0, len(info.State))
				for k := range info.State {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Printf("  %s: %v\n", k, info.State[k])
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
