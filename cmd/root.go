package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "nightel-core",
	Short: "The Nightel core!",
}

var banner = `Nightel Core`

func init() {
	cobra.OnInitialize(func() {
		fmt.Println(banner)
		fmt.Println(runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	})

	rootCmd.PersistentFlags().StringVarP(
		&configPath, "config", "c", "configs/config.yaml", "Path to configuration file",
	)

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
