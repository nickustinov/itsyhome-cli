package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show home status summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(config.Load())
		status, err := c.GetStatus()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(status, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Rooms:       %d\n", status.Rooms)
		fmt.Printf("Devices:     %d\n", status.Devices)
		fmt.Printf("Accessories: %d\n", status.Accessories)
		fmt.Printf("Reachable:   %d\n", status.Reachable)
		fmt.Printf("Unreachable: %d\n", status.Unreachable)
		fmt.Printf("Scenes:      %d\n", status.Scenes)
		fmt.Printf("Groups:      %d\n", status.Groups)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
