package cmd

import (
	"fmt"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nightel-core",
	Short: "The Nightel core!",
}

func init() {
	cobra.OnInitialize(func() { fmt.Println(config.AppName) })
}

func Execute() error {
	return rootCmd.Execute()
}
