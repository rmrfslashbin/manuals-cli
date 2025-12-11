package cmd

import (
	"fmt"

	"github.com/rmrfslashbin/manuals-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	devicesLimit  int
	devicesOffset int
	devicesDomain string
	devicesType   string
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List and manage devices",
	Long:  `List, view, and manage device documentation in the Manuals database.`,
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all devices",
	Long: `List devices in the Manuals database with optional filtering.

Filter by domain (hardware, software) or type (dev-boards, sensors, etc.).`,
	Example: `  manuals devices list
  manuals devices list --domain hardware
  manuals devices list --type dev-boards --limit 10
  manuals devices list -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ListDevices(devicesLimit, devicesOffset, devicesDomain, devicesType)
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}

		if out.IsJSON() {
			return out.JSON(result)
		}

		if len(result.Data) == 0 {
			out.Println("No devices found.")
			return nil
		}

		out.Text("Showing %d of %d devices:\n\n", len(result.Data), result.Total)

		headers := []string{"ID", "NAME", "DOMAIN", "TYPE"}
		rows := make([][]string, len(result.Data))
		for i, d := range result.Data {
			rows[i] = []string{
				d.ID[:8],
				output.Truncate(d.Name, 45),
				d.Domain,
				d.Type,
			}
		}
		out.Table(headers, rows)

		if result.Total > len(result.Data) {
			out.Text("\nUse --offset %d to see more results.\n", result.Offset+len(result.Data))
		}

		return nil
	},
}

var devicesGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get device details",
	Long:  `Get detailed information about a specific device by ID.`,
	Example: `  manuals devices get abc12345
  manuals devices get abc12345 -o json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		device, err := apiClient.GetDevice(id)
		if err != nil {
			return fmt.Errorf("failed to get device: %w", err)
		}

		if out.IsJSON() {
			return out.JSON(device)
		}

		out.Text("Device: %s\n", device.Name)
		out.Text("  ID:        %s\n", device.ID)
		out.Text("  Domain:    %s\n", device.Domain)
		out.Text("  Type:      %s\n", device.Type)
		out.Text("  Path:      %s\n", device.Path)
		out.Text("  Indexed:   %s\n", device.IndexedAt)

		if device.Content != "" {
			out.Text("\n--- Content ---\n%s\n", device.Content)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesGetCmd)

	devicesListCmd.Flags().IntVarP(&devicesLimit, "limit", "l", 50, "maximum number of results")
	devicesListCmd.Flags().IntVar(&devicesOffset, "offset", 0, "offset for pagination")
	devicesListCmd.Flags().StringVarP(&devicesDomain, "domain", "d", "", "filter by domain (hardware, software)")
	devicesListCmd.Flags().StringVarP(&devicesType, "type", "t", "", "filter by type")
}
