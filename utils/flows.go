/**
* We need a domain object defind in Go for every object we're using in the plan.
* There's always something that needs to be mutated. The domain object helps to do that.
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

// Flow is the Flow Domain Object
type Flow struct {
	ID          string      `json:"id,omitempty"`
	Path        string      `json:"path,omitempty"`
	Filter      string      `json:"filter,omitempty"`
	Description string      `json:"description,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Name        string      `json:"name,omitempty"`
	Capacity    int         `json:"capacity,omitempty"`
	Public      bool        `json:"public,omitempty"`
}

// RemoveFlowIDs removes the ids from the flows
func RemoveFlowIDs() {
	wd, _ := os.Getwd()

	filepath.Walk(wd, removeIDs)
}

func removeIDs(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	firstExt := filepath.Ext(info.Name())
	secondExt := filepath.Ext(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))
	if firstExt == ".json" && secondExt == ".flow" {
		flow := GetFlow(path)

		flow.ID = ""
		writeFlowJSON(flow.Name, flow)
	}

	return nil
}

func writeFlowJSON(filename string, flow Flow) {
	filename = filename + ".json"

	flow.Name = ""

	b, err := json.MarshalIndent(flow, "", "    ")
	Checkerror(err)

	if file, err := os.Create(filename); err != nil && !os.IsNotExist(err) {
		Checkerror(err)
	} else {
		defer file.Close()
		file.Write(b)
	}
}

// GetFlow gets the flow
func GetFlow(name string) Flow {
	data, err := ioutil.ReadFile(name)
	Checkerror(err)

	var flow Flow
	err = json.Unmarshal(data, &flow)
	Checkerror(err)

	barename := strings.TrimSuffix(name, filepath.Ext(name))

	flow.Name = barename

	return flow
}
