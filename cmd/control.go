package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nickustinov/itsyhome-cli/internal/client"
	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/spf13/cobra"
)

func makeControlCmd(action, short string, minArgs int) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s <target>", action),
		Short: short,
		Args:  cobra.MinimumNArgs(minArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := strings.Join(args, " ")
			path := "/" + action + "/" + target
			return doControl(path)
		},
	}
}

func makeValueControlCmd(action, short, valueDesc string) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s <%s> <target>", action, valueDesc),
		Short: short,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			value := args[0]
			target := strings.Join(args[1:], " ")
			path := "/" + action + "/" + value + "/" + target
			return doControl(path)
		},
	}
}

func doControl(path string) error {
	c := client.New(config.Load())
	resp, err := c.DoAction(path)
	if err != nil {
		return err
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Println(resp.Status)
	return nil
}

func init() {
	rootCmd.AddCommand(makeControlCmd("toggle", "Toggle a device or group", 1))
	rootCmd.AddCommand(makeControlCmd("on", "Turn on a device or group", 1))
	rootCmd.AddCommand(makeControlCmd("off", "Turn off a device or group", 1))
	rootCmd.AddCommand(makeControlCmd("lock", "Lock a device", 1))
	rootCmd.AddCommand(makeControlCmd("unlock", "Unlock a device", 1))
	rootCmd.AddCommand(makeControlCmd("open", "Open a device (blinds, garage)", 1))
	rootCmd.AddCommand(makeControlCmd("close", "Close a device (blinds, garage)", 1))
	rootCmd.AddCommand(makeControlCmd("scene", "Activate a scene", 1))

	rootCmd.AddCommand(makeValueControlCmd("brightness", "Set brightness (0-100)", "value"))
	rootCmd.AddCommand(makeValueControlCmd("position", "Set position (0-100)", "value"))
	rootCmd.AddCommand(makeValueControlCmd("temp", "Set color temperature (140-500 mireds)", "value"))
	rootCmd.AddCommand(makeValueControlCmd("color", "Set color (hex)", "hex"))
}
