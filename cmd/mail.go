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
	defer a.Close()
	if err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
	if err = a.Init(); err != nil {
		panic(fmt.Sprintf("%+v\n", err))
	}
	fmt.Printf("%+v\n", a.Config)
	a.Mailer.Send("realmiladrahimi@gmail.com", "Hello from the other side!")
	fmt.Println("Done!")
}
