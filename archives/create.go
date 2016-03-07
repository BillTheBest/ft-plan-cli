/**
* This will create the archive that gets sent up to flowthings.
**/

package archives

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/jhoonb/archivex"
	"github.com/flowthings/ft-plan-cli/utils"
)

var tarBall *archivex.TarFile
var rootDir string

// CreateArchive creates the tar archive
func CreateArchive(tempDir string, planyml []byte) (tarFilePath string) {
	tarFilePath = MakeTarFilePath(tempDir)
	createTarFile(tarFilePath, planyml)
	return tarFilePath
}

func createTarFile(tarFilePath string, planyml []byte) {
	// it's very unlikely we get an error here.
	wd, _ := os.Getwd()
	rootDir = wd

	tarBall = new(archivex.TarFile)
	tarBall.Create(tarFilePath)
	defer closeTarFile(tarBall)

	filepath.Walk(wd, addFiles)

	tarBall.Add("plan.yml", planyml)
}

func addFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	firstExt := filepath.Ext(info.Name())
	secondExt := filepath.Ext(strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))

	if firstExt == ".json" {
		dir, name := filepath.Split(path)
		subDir := strings.Replace(dir, rootDir, "", 1)
		entryName := filepath.Join(subDir, name)

		if secondExt == ".flow" {
			err := tarBall.AddFileWithName(path, entryName)
			if err != nil {
				fmt.Println(err)
			}
		} else if secondExt == ".track" {
			track := utils.GetTrack(path)

			addTrack(entryName, track, tarBall)
		}
	}

	return nil
}

func addTrack(filename string, track utils.Track, tar *archivex.TarFile) {
	track.Name = ""

	b, err := json.Marshal(track)
	utils.Checkerror(err)

	err = tar.Add(filename, b)

	if err != nil {
		fmt.Println(err)
	}
}
