package cmd

import (
	"github.com/neatplex/nightel-core/internal/app"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

	if err = a.Boot(); err != nil {
		a.Logger.Engine.Fatal("cannot start the app", zap.Error(err))
	}

	a.HttpServer.Serve()
	a.Wait()
}
