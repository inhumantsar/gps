package main

import (
	"fmt"
	"gps/internal/config"
	"gps/internal/gpt"
	"os"

	"github.com/spf13/cobra"
)

var cfg *config.Config
var cfgPath string
var opts gpt.ProjectOptions

var rootCmd = &cobra.Command{
	Use:   "gps",
	Short: "GPS on the command line",
	Long:  `Create new projects interactively.`,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application",
	Long:  `Create a new application with the specified name and language.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		resp, err := gpt.NewProject(cfg.Gpt, &opts)
		if err != nil {
			return err
		}
		fmt.Println(resp.Preview)
		fmt.Printf("\n---\n\n")
		fmt.Println(resp.Files)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, fmt.Sprintf("config file (default is %s)", config.DefaultPath))
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name of the application (required)")
	createCmd.Flags().StringVarP(&opts.Prompt, "prompt", "p", "", "Optional prompt to use for project creation. If a path is provided, the prompt will be loaded from that file.")
	createCmd.MarkFlagRequired("name")
}

func loadConfig() error {
	var err error
	if cfgPath != "" {
		cfg, err = config.LoadConfig(cfgPath)
	} else {
		cfg, err = config.LoadDefaultConfig()
	}
	return err
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
