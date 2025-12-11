// Package output handles formatting CLI output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents an output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatText  Format = "text"
)

// Writer handles formatted output.
type Writer struct {
	format Format
	out    io.Writer
}

// New creates a new output writer.
func New(format string) *Writer {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatJSON, FormatText:
		// valid
	default:
		f = FormatTable
	}
	return &Writer{
		format: f,
		out:    os.Stdout,
	}
}

// JSON outputs data as JSON.
func (w *Writer) JSON(data interface{}) error {
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Table outputs data as a table.
func (w *Writer) Table(headers []string, rows [][]string) {
	tw := tabwriter.NewWriter(w.out, 0, 0, 2, ' ', 0)
	
	// Print headers
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	fmt.Fprintln(tw, strings.Repeat("-", len(strings.Join(headers, "  "))))
	
	// Print rows
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	
	tw.Flush()
}

// Text outputs plain text.
func (w *Writer) Text(format string, args ...interface{}) {
	fmt.Fprintf(w.out, format, args...)
}

// Println outputs a line.
func (w *Writer) Println(args ...interface{}) {
	fmt.Fprintln(w.out, args...)
}

// Format returns the current output format.
func (w *Writer) Format() Format {
	return w.format
}

// IsJSON returns true if the output format is JSON.
func (w *Writer) IsJSON() bool {
	return w.format == FormatJSON
}

// Truncate truncates a string to a maximum length.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// FormatSize formats a byte size as a human-readable string.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
