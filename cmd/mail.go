package cmd

import (
	"fmt"

	"github.com/neatplex/nightell-core/internal/app"
	"github.com/spf13/cobra"
)

var mailCmd = &cobra.Command{
	Use:   "mail [to]",
	Short: "Send mail through the application for testing purpose",
	Args:  cobra.MaximumNArgs(1),
	Run:   mailFunc,
}

func init() {
	rootCmd.AddCommand(mailCmd)
}

func mailFunc(_ *cobra.Command, args []string) {
	to := "realmiladrahimi@gmail.com"
	if len(args) > 0 {
		to = args[0]
	}

	a, err := app.New()
	if err != nil {
		if a != nil {
			a.Close()
		}
		panic(fmt.Sprintf("%+v\n", err))
	}
	defer a.Close()

	a.Container.Mailer.Send(to, "Hello", "Hello from the other side!")
	fmt.Println("Done!")
}
