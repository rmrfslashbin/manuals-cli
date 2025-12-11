package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rmrfslashbin/manuals-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	docsLimit    int
	docsOffset   int
	docsDeviceID string
	docsOutput   string
)

var documentsCmd = &cobra.Command{
	Use:     "documents",
	Aliases: []string{"docs"},
	Short:   "List and download documents",
	Long:    `List and download documentation files (PDFs, datasheets, etc.).`,
}

var documentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all documents",
	Long: `List documents in the Manuals database.

Filter by device ID to see documents for a specific device.`,
	Example: `  manuals documents list
  manuals docs list --device abc12345
  manuals docs list --limit 20 -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ListDocuments(docsLimit, docsOffset, docsDeviceID)
		if err != nil {
			return fmt.Errorf("failed to list documents: %w", err)
		}

		if out.IsJSON() {
			return out.JSON(result)
		}

		if len(result.Data) == 0 {
			out.Println("No documents found.")
			return nil
		}

		out.Text("Showing %d of %d documents:\n\n", len(result.Data), result.Total)

		headers := []string{"ID", "FILENAME", "TYPE", "SIZE"}
		rows := make([][]string, len(result.Data))
		for i, d := range result.Data {
			rows[i] = []string{
				d.ID[:8],
				output.Truncate(d.Filename, 45),
				d.MimeType,
				output.FormatSize(d.SizeBytes),
			}
		}
		out.Table(headers, rows)

		if result.Total > len(result.Data) {
			out.Text("\nUse --offset %d to see more results.\n", result.Offset+len(result.Data))
		}

		return nil
	},
}

var documentsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get document details",
	Long:  `Get detailed information about a specific document by ID.`,
	Example: `  manuals docs get abc12345
  manuals documents get abc12345 -o json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		doc, err := apiClient.GetDocument(id)
		if err != nil {
			return fmt.Errorf("failed to get document: %w", err)
		}

		if out.IsJSON() {
			return out.JSON(doc)
		}

		out.Text("Document: %s\n", doc.Filename)
		out.Text("  ID:        %s\n", doc.ID)
		out.Text("  Device:    %s\n", doc.DeviceID)
		out.Text("  Path:      %s\n", doc.Path)
		out.Text("  Type:      %s\n", doc.MimeType)
		out.Text("  Size:      %s\n", output.FormatSize(doc.SizeBytes))
		out.Text("  Checksum:  %s\n", doc.Checksum[:16]+"...")
		out.Text("  Indexed:   %s\n", doc.IndexedAt)

		return nil
	},
}

var documentsDownloadCmd = &cobra.Command{
	Use:   "download <id>",
	Short: "Download a document",
	Long: `Download a document file by ID.

By default, saves to the current directory with the original filename.
Use --output to specify a different path.`,
	Example: `  manuals docs download abc12345
  manuals docs download abc12345 -o ~/Documents/datasheet.pdf
  manuals documents download abc12345 --output ./docs/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		// Get document info first for the filename
		doc, err := apiClient.GetDocument(id)
		if err != nil {
			return fmt.Errorf("failed to get document info: %w", err)
		}

		// Download the file
		body, _, err := apiClient.DownloadDocument(id)
		if err != nil {
			return fmt.Errorf("failed to download document: %w", err)
		}
		defer body.Close()

		// Determine output path
		outputPath := docsOutput
		if outputPath == "" {
			outputPath = doc.Filename
		} else {
			// Check if output is a directory
			info, err := os.Stat(outputPath)
			if err == nil && info.IsDir() {
				outputPath = filepath.Join(outputPath, doc.Filename)
			}
		}

		// Create output file
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()

		// Copy content
		written, err := io.Copy(file, body)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		out.Text("Downloaded %s (%s) to %s\n", doc.Filename, output.FormatSize(written), outputPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(documentsCmd)
	documentsCmd.AddCommand(documentsListCmd)
	documentsCmd.AddCommand(documentsGetCmd)
	documentsCmd.AddCommand(documentsDownloadCmd)

	documentsListCmd.Flags().IntVarP(&docsLimit, "limit", "l", 50, "maximum number of results")
	documentsListCmd.Flags().IntVar(&docsOffset, "offset", 0, "offset for pagination")
	documentsListCmd.Flags().StringVar(&docsDeviceID, "device", "", "filter by device ID")

	documentsDownloadCmd.Flags().StringVarP(&docsOutput, "output", "o", "", "output path (file or directory)")
}
