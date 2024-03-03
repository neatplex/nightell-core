package cmd

import (
	"fmt"
	"github.com/neatplex/nightell-core/internal/app"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the http server.",
	Run:   serveFunc,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serveFunc(_ *cobra.Command, _ []string) {
	a, err := app.New()
	defer a.Close()
	if err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
	if err = a.Init(); err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
	a.HttpServer.Serve()
	a.Wait()
}
