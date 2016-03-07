/**
* We need a domain object defined in Go for every object we're using in the plan.
* There's always something that needs to be mutated. The domain object helps to do that.
* In this case, we mostly want to mutate the JavaScript.
* We want to separate it from the reset of the JSON blob and make it its own file.
**/

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Track is the Track Domain Object
type Track struct {
	ID          string      `json:"id,omitempty"`
	Source      string      `json:"source,omitempty"`
	Destination string      `json:"destination,omitempty"`
	Filter      string      `json:"filter,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Js          string      `json:"js,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
}

// AddJavaScript adds javascript to the tracks
func AddJavaScript() {
	wd, _ := os.Getwd()
	trackDirPath := filepath.Join(wd, "flows")

	filepath.Walk(trackDirPath, addJavascript)
}

// SeparateJavaScript separates javascript from the tracks
func SeparateJavaScript() {
	wd, _ := os.Getwd()
	trackDirPath := filepath.Join(wd, "flows")

	filepath.Walk(trackDirPath, separateJavascript)
}

func addJavascript(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	firstExt := filepath.Ext(info.Name())
	secondExt := filepath.Ext(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))
	if firstExt == ".json" && secondExt == ".track" {
		track := GetTrack(path)
		jsname := track.Name + ".js"
		os.Remove(jsname)

		writeTrackJSON(track.Name, track)
	}

	return nil
}

func separateJavascript(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	firstExt := filepath.Ext(info.Name())
	secondExt := filepath.Ext(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))
	if firstExt == ".json" && secondExt == ".track" {
		track := GetTrack(path)
		js := track.Js

		writeJavaScript(track.Name, js)
		track.Js = ""
		track.ID = ""
		writeTrackJSON(track.Name, track)
	}

	return nil
}

// GetTrack gets the track
func GetTrack(name string) Track {
	data, err := ioutil.ReadFile(name)
	Checkerror(err)

	var track Track
	err = json.Unmarshal(data, &track)
	Checkerror(err)

	barename := strings.TrimSuffix(name, filepath.Ext(name))

	track.Name = barename

	jsname := barename + ".js"

	if FileExists(jsname) {
		jsdata, err := ioutil.ReadFile(jsname)
		Checkerror(err)

		js := string(jsdata[:])
		track.Js = js
	}

	return track
}

func writeTrackJSON(filename string, track Track) {
	filename = filename + ".json"

	track.Name = ""

	b, err := json.MarshalIndent(track, "", "    ")
	Checkerror(err)

	if file, err := os.Create(filename); err != nil && !os.IsNotExist(err) {
		Checkerror(err)
	} else {
		defer file.Close()
		file.Write(b)
	}
}

func writeJavaScript(filename string, js string) {
	if len(js) == 0 {
		return
	}
	filename = filename + ".js"
	if file, err := os.Create(filename); err != nil && !os.IsNotExist(err) {
		Checkerror(err)
	} else {
		defer file.Close()
		jsB := []byte(js)
		file.Write(jsB)
	}
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
