package cmd

import (
	"log"

	"github.com/samuelngs/wordpress-cli/app"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View output from containers.",
	Run: func(cmd *cobra.Command, args []string) {

		app, err := app.New()
		if err != nil {
			log.Fatal(err)
		}

		switch follow, err := cmd.Flags().GetBool("follow"); {
		case err != nil:
			log.Fatal(err)
		default:
			if err := app.Log(follow); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output.")

}
