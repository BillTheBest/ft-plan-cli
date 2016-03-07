package flowthings

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/flowthings/ft-plan-cli/utils"
)

func addHeaders(request *http.Request) *http.Request {
	request.Header.Add("X-Auth-Account", config.username)
	request.Header.Add("X-Auth-Token", config.token)
	request.Header.Add("User-Agent", userAgent)
	return request
}

func createTempDirectory(dir string, inDev bool) string {
	var tempDir string
	if inDev {
		tempDir = path.Join(dir, "tmp")
	} else {
		tempDir = path.Join(os.TempDir(), "ft")
	}

	if err := os.MkdirAll(tempDir, 0777); err != nil && !os.IsExist(err) {
		utils.Checkerror(err)
	}
	return tempDir
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
