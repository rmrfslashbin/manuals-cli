package cmd

import (
	"fmt"
	"strings"

	"github.com/rmrfslashbin/manuals-cli/internal/output"
	"github.com/spf13/cobra"
)

var searchLimit int

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for devices and documentation",
	Long: `Search the Manuals database for devices matching your query.

Uses semantic (vector) search to find relevant hardware and software documentation.
Results are ranked by relevance and include snippet previews.`,
	Example: `  manuals search "raspberry pi gpio"
  manuals search "uart protocol" --limit 5
  manuals search esp32 -o json`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		results, err := apiClient.Search(query, searchLimit)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if out.IsJSON() {
			return out.JSON(results)
		}

		if len(results.Results) == 0 {
			out.Println("No results found.")
			return nil
		}

		out.Text("Found %d results for \"%s\":\n\n", results.Total, results.Query)

		headers := []string{"ID", "NAME", "DOMAIN", "TYPE", "SCORE"}
		rows := make([][]string, len(results.Results))
		for i, r := range results.Results {
			rows[i] = []string{
				r.DeviceID[:8],
				output.Truncate(r.Name, 40),
				r.Domain,
				r.Type,
				fmt.Sprintf("%.2f", r.Score),
			}
		}
		out.Table(headers, rows)

		// Show snippets for top results
		if len(results.Results) > 0 && outputFormat != "table" {
			out.Println("\n--- Snippets ---")
			for i, r := range results.Results {
				if i >= 3 {
					break
				}
				if r.Snippet != "" {
					out.Text("\n[%s] %s\n", r.DeviceID[:8], r.Name)
					out.Text("  %s\n", output.Truncate(r.Snippet, 200))
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 20, "maximum number of results")
}
