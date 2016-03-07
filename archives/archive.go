/**
* Archive utils contains the functions to pack and unpack the archive.
* It also have all of the file manipulation logic in it.
**/

package archives

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/jhoonb/archivex"
	"github.com/flowthings/ft-plan-cli/utils"
)

func closeTarFile(tar *archivex.TarFile) {
	err := tar.Close()
	utils.Checkerror(err)
}

// "/a/b" -> /a/b/plans.xys123.tar.gz"
func MakeTarFilePath(tempDir string) string {
	return filepath.Join(tempDir, makeUniqueFileName("plans", "tar.gz"))
}

//"a","b" -> "a.12823ajsh.b
func makeUniqueFileName(prefix string, suffix string) string {
	return strings.Join([]string{prefix, strconv.FormatInt(int64(time.Now().Unix()), 32), suffix}, ".")
}
