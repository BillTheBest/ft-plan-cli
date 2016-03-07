package flowthings

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/flowthings/ft-plan-cli/archives"
	"github.com/flowthings/ft-plan-cli/utils"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull down updates to your plan from flowthings",
	Long: `If you edit your project on flowthings,

you'll want to be able to pull down a local copy.
This command pulls down a copy to your local machine and syncs the changes for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkPlanExists(config)

		var archivePath string = retrieveArchive()
		planConfig := archives.UnpackArchive(archivePath, config.tempDir)

		id := planConfig["plan_id"]

		message := "Pulled plan " + id + " from " + config.basePath

		color.Green(message)
	},
}

func retrieveArchive() (archivePath string) {
	client := &http.Client{}
	request := newGetArchiveRequest()

	resp, err := client.Do(request)
	utils.Checkerror(err)
	defer resp.Body.Close()

	contents := &bytes.Buffer{}
	_, err = contents.ReadFrom(resp.Body)
	utils.Checkerror(err)

	utils.BadStatusCode(resp, contents)

	var dat map[string]interface{}
	err = json.Unmarshal(contents.Bytes(), &dat)
	if err == nil {
		// we only want to do something if it parsed,
		// because otherwise it's probably not a JSON object.
		var head = dat["head"].(map[string]interface{})
		var errs = head["errors"].([]interface{})
		for _, err := range errs {
			color.Red(err.(string))
		}
		os.Exit(1)
	}

	archivePath = archives.MakeTarFilePath(config.tempDir)

	file, err := os.Create(archivePath)
	defer file.Close()

	file.Write(contents.Bytes())
	return archivePath
}

func newGetArchiveRequest() (request *http.Request) {
	rawURLStr := planURI(config) + "?path=" + url.QueryEscape(config.basePath)
	request, err := http.NewRequest("PUT", rawURLStr, nil)
	utils.Checkerror(err)
	request = addHeaders(request)
	return request
}

// checkPlanExists will check to see if the user has provided a base path
func checkPlanExists(c *configuration) {
	if c.basePath == "" {
		color.Red("You cannot pull without a basepath")
		os.Exit(1)
	}
}
