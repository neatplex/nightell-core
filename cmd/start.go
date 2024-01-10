package cmd

import (
	"github.com/neatplex/nightel-core/internal/app"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the service.",
	Run:   startFunc,
}

func startFunc(_ *cobra.Command, _ []string) {
	a, err := app.New(configPath)
	if err != nil {
		panic(err)
	}

	a.Boot()
	a.Wait()
}
