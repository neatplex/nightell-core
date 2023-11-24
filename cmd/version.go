package cmd

import (
	"fmt"
	"github.com/neatplex/nightel-core/internal/config"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version of the application",
	Run:   versionFunc,
}

func versionFunc(_ *cobra.Command, _ []string) {
	fmt.Println(config.AppVersion, "[", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH, "]")
}
