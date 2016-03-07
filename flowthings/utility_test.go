package flowthings

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestCreateDevTempDir(t *testing.T) {
	os.Remove(".test")
	os.MkdirAll(".test", 0777)
	cwd, _ := os.Getwd()
	createTempDirectory(".test/tmp", true)
	if _, err := os.Stat(".test/tmp"); os.IsNotExist(err) {
		t.Errorf("Fail creating dev temp dir %s", path.Join(cwd, ".test/tmp"))
	}
	os.Remove(".test")

}
func TestCreateStdTempDir(t *testing.T) {
	dir := createTempDirectory("this should be ignored", false)
	assert.True(t, fileExists(dir), fmt.Sprintf("temp dir %s exists", dir))
}
