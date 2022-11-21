package cmd

import (
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/config"
	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/hubble-middleware/cmd/api"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new iot-platform root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "hubble-middleware",
	}

	cfg := config.New()

	api.Register(root, cfg)

	return root
}
