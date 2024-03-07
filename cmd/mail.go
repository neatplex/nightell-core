package cmd

import (
	"github.com/neatplex/nightel-core/internal/mailer"
	"github.com/spf13/cobra"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Send mail through the application",
	Run:   mailFunc,
}

func init() {
	rootCmd.AddCommand(mailCmd)
}

func mailFunc(_ *cobra.Command, _ []string) {
	mailer.Send()
}
