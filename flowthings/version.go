package flowthings

import (
	"fmt"
	"os"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of flowthings",
	Long:  "Print the version of flowthings plan CLI",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("Flowthings Plan CLI Version %s\n", ftCliVersion)
	os.Exit(0)

}
