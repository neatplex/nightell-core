package cmd

import (
	"fmt"

	"github.com/neatplex/nightell-core/internal/app"
	"github.com/spf13/cobra"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Send mail through the application for testing purpose",
	Run:   mailFunc,
}

func init() {
	rootCmd.AddCommand(mailCmd)
}

func mailFunc(_ *cobra.Command, _ []string) {
	a, err := app.New()
	if err != nil {
		if a != nil {
			a.Close()
		}
		panic(fmt.Sprintf("%+v\n", err))
	}
	defer a.Close()

	fmt.Printf("%+v\n", a.Container.Config)
	a.Container.Mailer.Send("realmiladrahimi@gmail.com", "Hello", "Hello from the other side!")
	fmt.Println("Done!")
}
