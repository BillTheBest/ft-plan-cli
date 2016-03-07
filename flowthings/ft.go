package flowthings

import (
	"fmt"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

var ftCliVersion = "1.0"
var defaultPlatformVersion = "v0.1"
var defaultPlatformEndpoint = "https://api.flowthings.io"
var userAgent = "Go-Plan-CLI"
var config = defaultConfig()
var inputValues = configuration{}

//------------------------------------------
// Main Command setup
//------------------------------------------

var ftCmd = &cobra.Command{
	Use:   "ft",
	Short: "flowthings plan installer",
	Long:  `ft is the main command. It helps to build and manage your flowthings Plans. Complete documentation is available at flowthings.io/docs`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		updateConfig(&inputValues, config)
		if err := errIfInvalid(config); err != nil {
			color.Red(fmt.Sprintf("Error: %s", err.Error()))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ftCmd.PersistentFlags().StringVarP(&inputValues.username, "username", "u", "", "Your username")
	ftCmd.PersistentFlags().StringVarP(&inputValues.token, "token", "t", "", "Your token")
	ftCmd.PersistentFlags().StringVarP(&inputValues.ftEndpoint, "endpoint", "e", defaultPlatformEndpoint, "Api Endpoint")
	ftCmd.PersistentFlags().StringVarP(&inputValues.ftVersion, "api-version", "v", defaultPlatformVersion, "Api Version")
	ftCmd.PersistentFlags().BoolVarP(&inputValues.ftClientDev, "dev", "x", false, "CLI Developement Version")
	ftCmd.PersistentFlags().StringVarP(&inputValues.configFile, "config", "c", config.configFile, "Config file")
	ftCmd.PersistentFlags().StringVarP(&inputValues.basePath, "basepath", "p", "", "Plan Base Path")
	ftCmd.PersistentFlags().StringVarP(&inputValues.planDir, "directory", "d", config.planDir, "Plan Directory")
	ftCmd.AddCommand(versionCmd)

	ftCmd.AddCommand(loginCmd)
	ftCmd.AddCommand(pushCmd)
	ftCmd.AddCommand(pullCmd)
}

// Execute is the ft entry point. This is called by the CLI main
func Execute(version string) {
	// This just proxies to the ftCmd.
	// We're using it because it allows the internal api to change without affecting the exported API.
	if version != "" {
		ftCliVersion = version
	}
	ftCmd.Execute()
}
