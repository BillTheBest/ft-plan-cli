package flowthings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/fatih/color"
	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/flowthings/ft-plan-cli/utils"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/flowthings/ft-plan-cli/archives"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push up to flowthings",
	Long:  "Push your plan up to the flowthings platform",
	Run: func(cmd *cobra.Command, args []string) {
		plan := baseplanYaml()
		var tarFilePath string = archives.CreateArchive(config.tempDir, plan)

		tarFile, err := os.Open(tarFilePath)
		utils.Checkerror(err)
		defer utils.CloseAndDelete(tarFile, tarFilePath)
		request := newSendArchiveRequest(tarFile, tarFilePath)

		sendArchive(request)
	},
}

// we don't want to add the base level config file raw
// we want to add the parsed yaml file, so that we can replace
// any of the relevant bits.
func baseplanYaml() []byte {
	plan := map[string]string{
		"base_path": config.basePath,
	} //

	planyaml, err := yaml.Marshal(&plan)
	utils.Checkerror(err)

	return planyaml
}

func sendArchive(request *http.Request) {
	client := &http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		utils.Checkerror(err)

		utils.BadStatusCode(resp, body)

		resp.Body.Close()
		// fmt.Println(resp.StatusCode)
		// fmt.Println(resp.Header)
		// fmt.Println(body)

		var dat map[string]interface{}
		err = json.Unmarshal(body.Bytes(), &dat)
		if err == nil {
			// we only want to do something if it parsed,
			// because otherwise it's probably not a JSON object.
			var body = dat["body"].(map[string]interface{})

			var id = body["id"].(string)
			origplanconfig := make(map[string]string)
			wd, _ := os.Getwd()
			byarr, err := ioutil.ReadFile(path.Join(wd, "plan.yml"))
			if err != nil && !os.IsNotExist(err) {
				utils.Checkerror(err)
			}
			err = yaml.Unmarshal(byarr, origplanconfig)
			origplanconfig["plan_id"] = id

			yml, err := yaml.Marshal(origplanconfig)
			utils.Checkerror(err)

			file, err := os.Create(path.Join(wd, "plan.yml"))
			utils.Checkerror(err)

			if _, err := file.Write(yml); err != nil {
				utils.Checkerror(err)
			}

			message := "Pushed plan " + id + " to " + config.basePath
			color.Green(message)
		}
	}
}

func newSendArchiveRequest(tarFile *os.File, tarFilePath string) (request *http.Request) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("filedata", filepath.Base(tarFilePath))

	_, err = io.Copy(part, tarFile)
	utils.Checkerror(err)

	contentType := writer.FormDataContentType()

	err = writer.Close()
	utils.Checkerror(err)

	request, err = http.NewRequest("POST", planURI(config), body)
	utils.Checkerror(err)

	request = addHeaders(request)
	request.Header.Add("Content-Type", contentType)
	return request
}
