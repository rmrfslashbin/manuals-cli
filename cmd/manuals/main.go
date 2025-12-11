// manuals is a CLI tool for interacting with the Manuals documentation platform.
package main

import (
	"fmt"
	"os"

	"github.com/rmrfslashbin/manuals-cli/cmd/manuals/cmd"
)

// Build information set via ldflags.
var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, gitCommit, buildTime)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
