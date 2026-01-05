package ritual

import (
	ritualcli "github.com/toutaio/toutago-ritual-grove/pkg/cli"
	"github.com/spf13/cobra"
)

// RitualCommand returns the ritual command from ritual-grove
func RitualCommand() *cobra.Command {
	return ritualcli.RitualCommand()
}
