package cmd

import (
	"log"

	"github.com/samuelngs/wordpress-cli/app"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove stopped containers",
	Long: `Removes stopped service containers.

By default, anonymous volumes attached to containers will not be removed. You
can override this with -v. To list all volumes, use docker volume ls.

Any data which is not in a volume will be lost.`,
	Run: func(cmd *cobra.Command, args []string) {

		app, err := app.New()
		if err != nil {
			log.Fatal(err)
		}

		if err := app.Remove(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
