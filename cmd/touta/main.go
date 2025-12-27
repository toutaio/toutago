package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toutaio/toutago/internal/cli"
)

var version = "0.1.0"

func main() {
	root := &cobra.Command{
		Use:   "touta",
		Short: "Toutā - A message-driven Go web framework",
		Long: `Toutā is a Go web framework that emphasizes:
  - Interface-first design for pluggability
  - Message-passing architecture
  - Dependency injection`,
		Version: version,
	}

	// Add commands
	root.AddCommand(cli.NewCommand())
	root.AddCommand(cli.InitCommand())
	root.AddCommand(cli.ServeCommand())
	root.AddCommand(cli.VersionCommand(version))

	// TODO: Dynamically load additional commands from plugins

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
