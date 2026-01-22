package cmd

import (
	"fmt"

	"github.com/nickustinov/itsyhome-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or update CLI configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		fmt.Printf("Host: %s\n", cfg.Host)
		fmt.Printf("Port: %d\n", cfg.Port)
		fmt.Printf("URL:  %s\n", cfg.BaseURL())
		fmt.Printf("File: %s\n", config.Path())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		if host, _ := cmd.Flags().GetString("host"); host != "" {
			cfg.Host = host
		}
		if port, _ := cmd.Flags().GetInt("port"); port != 0 {
			cfg.Port = port
		}

		if err := config.Save(cfg); err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println("Configuration saved.")
	},
}

func init() {
	configSetCmd.Flags().String("host", "", "Server host address")
	configSetCmd.Flags().Int("port", 0, "Server port")
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
