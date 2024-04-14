package cmd

import (
	"fmt"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version of the application",
	Run:   versionFunc,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionFunc(_ *cobra.Command, _ []string) {
	fmt.Println(config.AppVersion, "(", runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH, ")")
}
