package cmd

import (
	"fmt"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nightell-core",
	Short: "The Nightel core!",
}

func init() {
	cobra.OnInitialize(func() { fmt.Println(config.AppName) })
}

func Execute() error {
	return rootCmd.Execute()
}
