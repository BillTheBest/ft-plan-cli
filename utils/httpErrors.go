/**
* This checks the http call to see if there were any errors in it.
**/

package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/fatih/color"
)

// TODO: We should provide more information on a failure.
// Perhaps we should print out the status code.

// BadStatusCode parses whether or not the response was an improper response.
// If it's not a valid response code, it will print the errors and exit.
func BadStatusCode(resp *http.Response, body *bytes.Buffer) {
	// This should check all failures.
	if resp.StatusCode >= 400 {
		var dat map[string]interface{}
		// we've gotta see what the platform said about the error.
		// so we parse the JSON it sent back in the body.
		err := json.Unmarshal(body.Bytes(), &dat)
		if err == nil {
			// we only want to do something with the JSON if it parsed without error.
			var head = dat["head"].(map[string]interface{})
			var errs = head["errors"].([]interface{})
			for _, err := range errs {
				color.Red(err.(string))
			}
			os.Exit(1)
		} else {
			// however, we still want to display an error, regardless
			color.Red(err.Error())
			os.Exit(1)
		}
	}
}
