// Package cmd implements the CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/rmrfslashbin/manuals-cli/internal/client"
	"github.com/rmrfslashbin/manuals-cli/internal/config"
	"github.com/rmrfslashbin/manuals-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	// Version information
	version   string
	gitCommit string
	buildTime string

	// Global flags
	cfgFile      string
	apiURL       string
	apiKey       string
	outputFormat string

	// Global state
	cfg    *config.Config
	apiClient *client.Client
	out    *output.Writer
)

// SetVersionInfo sets the version information.
func SetVersionInfo(v, commit, build string) {
	version = v
	gitCommit = commit
	buildTime = build
}

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "manuals",
	Short: "CLI for the Manuals documentation platform",
	Long: `manuals is a command-line interface for searching and accessing
hardware and software documentation from the Manuals platform.

Configure the API endpoint and key via environment variables:
  MANUALS_API_URL  - API base URL (default: http://localhost:8080)
  MANUALS_API_KEY  - API key for authentication (required)

Or create a config file at ~/.manuals.yaml:
  api_url: http://manuals.local:8080
  api_key: your-api-key
  output_format: table`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for version and help commands
		if cmd.Name() == "version" || cmd.Name() == "help" {
			return nil
		}

		// Load configuration
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}

		// Override with flags
		if apiURL != "" {
			cfg.APIBaseURL = apiURL
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}
		if outputFormat != "" {
			cfg.OutputFormat = outputFormat
		}

		// Validate
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Initialize client and output
		apiClient = client.New(cfg.APIBaseURL, cfg.APIKey)
		out = output.New(cfg.OutputFormat)

		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.manuals.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format (table, json, text)")
}

// versionCmd shows version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("manuals version %s\n", version)
		fmt.Printf("  commit: %s\n", gitCommit)
		fmt.Printf("  built:  %s\n", buildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// exitError prints an error and exits.
func exitError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
