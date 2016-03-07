package flowthings

import (
	"os"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login -u [username] -p [password]",
	Short: "store your login credentials",
	Long:  "Store your login credentials for the flowthings platform locally. At least a password and a token must be provided.",
	Run: func(cmd *cobra.Command, args []string) {
		if inputValues.username == "" && inputValues.token == "" {
			if !anythingToUpdate() {
				cmd.Help()
				os.Exit(0)
			}
		}
		login(config)
	},
}

func anythingToUpdate() bool {
	if inputValues.ftEndpoint != config.ftEndpoint {
		return true
	} else if inputValues.ftClientDev != config.ftClientDev {
		return true
	} else if inputValues.ftVersion != config.ftVersion {
		return true
	}
	return false
}

func login(c *configuration) {
	configToFile(c.configFile, c)
}
