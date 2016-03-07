/**
* This will unpack the archive that gets sent down from flowthings.
**/

package archives

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/flowthings/ft-plan-cli/utils"
)

// UnpackArchive unpacks the plan archive into the local filestructure
func UnpackArchive(tarPath string, tempDir string) map[string]string {
	removeExisting()
	planConfig := unpackArchive(tarPath, tempDir)

	utils.SeparateJavaScript()
	utils.RemoveFlowIDs()

	return planConfig
}

func unpackArchive(tarPath string, tempDir string) map[string]string {
	tarFile, err := os.Open(tarPath)
	utils.Checkerror(err)
	defer utils.CloseAndDelete(tarFile, tarPath)
	wd, _ := os.Getwd()
	rootDir = wd

	gzReader, err := gzip.NewReader(tarFile)
	utils.Checkerror(err)
	defer gzReader.Close()

	tarBallReader := tar.NewReader(gzReader)

	origplanconfig := make(map[string]string)

	for {
		// Get the next file in the tarball
		// if there isn't a next file,
		// break out of the loop.
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Checkerror(err)
		}

		filename := header.Name
		filename = path.Join(rootDir, filename)

		if header.Name != "plan.yml" {
			removeFile(filename)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(filename, os.FileMode(header.Mode))
			utils.Checkerror(err)
		case tar.TypeReg, tar.TypeRegA:
			dir := path.Dir(filename)
			err := os.MkdirAll(dir, os.FileMode(0777))
			utils.Checkerror(err)

			if header.Name == "plan.yml" {
				body := &bytes.Buffer{}
				body.ReadFrom(tarBallReader)
				var newplanconfig map[string]string
				yaml.Unmarshal(body.Bytes(), &newplanconfig)

				byarr, err := ioutil.ReadFile(path.Join(rootDir, "plan.yml"))
				if err != nil && !os.IsNotExist(err) {
					utils.Checkerror(err)
				}
				err = yaml.Unmarshal(byarr, origplanconfig)

				origplanconfig["base_path"] = newplanconfig["base_path"]
				origplanconfig["plan_id"] = newplanconfig["plan_id"]

				yml, err := yaml.Marshal(origplanconfig)
				utils.Checkerror(err)

				file, err := os.Create(path.Join(rootDir, "plan.yml"))
				utils.Checkerror(err)

				if _, err := file.Write(yml); err != nil {
					utils.Checkerror(err)
				}
			} else {
				writer, err := os.Create(filename)
				utils.Checkerror(err)

				io.Copy(writer, tarBallReader)

				err = os.Chmod(filename, os.FileMode(header.Mode))
				utils.Checkerror(err)

				writer.Close()
			}
		}
	}

	return origplanconfig
}

func removeExisting() {
	wd, _ := os.Getwd()
	// we don't want to remove the plan.yml,
	// because the new one won't be an exact match.

	// we just want to read that one in and (potentially)
	// merge it with the existing one.
	files := [...]string{"tracks", "flows", "tasks", "devices"}
	for _, file := range files {
		filePath := filepath.Join(wd, file)
		removeFile(filePath)
	}
}

/**
* We remove the existing files before we write the new ones to file.
* Diffing the files would be wasted effort, we're going to overwrite many of them, anyway.
**/
func removeFile(filename string) {
	fileInfo, _ := os.Stat(filename)
	// Does the file exist? If it's nil, it doesn't.
	if fileInfo == nil {
		return
	}

	// is the file a directory? Then we want to use the directory removal
	if fileInfo.IsDir() {
		if err := os.RemoveAll(filename); err != nil && !os.IsNotExist(err) {
			utils.Checkerror(err)
		}
		// if not, we'll use the normal file removal
	} else {
		if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
			utils.Checkerror(err)
		}
	}
}
