package main

import (
	"os"

	"gitlab.snapp.ir/snappcloud/hubble-middleware/internal/hubble-middleware/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
