package cmd

import (
	"fmt"
	"log"

	"github.com/samuelngs/wordpress-cli/app"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		app, err := app.New()
		if err != nil {
			log.Fatal(err)
		}

		if err := app.Up(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nServer is running on localhost:%s\n", app.Port())

		switch detach, err := cmd.Flags().GetBool("detach"); {
		case err != nil:
			log.Fatal(err)
		case !detach:
			if err := app.Log(true); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\ngracefully stopping...\n")
		}
	},
}

func init() {
	RootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolP("detach", "d", false, "Run containers in the background")
}
