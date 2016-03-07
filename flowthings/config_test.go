package flowthings

import (
	"os"
	"path"
	"testing"

	"os/exec"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var startingConfigFile string

func TestMain(m *testing.M) {

	//preserve existing creds if they exist
	defaultConfigFile := defaultConfigFile()
	if fileExists(defaultConfigFile) {
		startingConfigFile = path.Join(os.TempDir(), "saved_creds.yaml")
		copyFile(defaultConfigFile, startingConfigFile)
	}
	m.Run()
	// restore existing file if it exists.
	if fileExists(startingConfigFile) {
		copyFile(startingConfigFile, defaultConfigFile)
		os.Remove(startingConfigFile)
	}
}

func TestVerify(t *testing.T) {
	assert.Error(t, errIfInvalid(&configuration{username: "", token: "x"}))
	assert.Error(t, errIfInvalid(&configuration{username: "x", token: ""}))
	assert.Nil(t, errIfInvalid(&configuration{username: "x", token: "x"}))
}

func TestDefaultConfig(t *testing.T) {
	//with no default config file existsing:
	os.Remove(defaultConfigFile())
	dc := defaultConfig()
	assert.Equal(t, dc.ftEndpoint, defaultPlatformEndpoint)
	assert.Equal(t, dc.ftVersion, defaultPlatformVersion)
	assert.Equal(t, dc.ftClientDev, false)
	assert.Equal(t, dc.username, "")
	assert.Equal(t, dc.token, "")

	//with a created config file:
	cc := configuration{
		username:    "testuser",
		token:       "testtoken",
		ftEndpoint:  "testendpoint",
		ftClientDev: true,
		ftVersion:   "testversion",
	}
	//this will write a new .ft/creds.yaml file
	configToFile(defaultConfigFile(), &cc)
	dc = defaultConfig()
	assert.Equal(t, dc.ftEndpoint, "testendpoint")
	assert.Equal(t, dc.ftVersion, "testversion")
	assert.Equal(t, dc.ftClientDev, true)
	assert.Equal(t, dc.username, "testuser")
	assert.Equal(t, dc.token, "testtoken")
}

func TestSpecifiedConfig(t *testing.T) {
	testConfigFile := path.Join(os.TempDir(), "test_creds.yaml")
	//with a created config file:
	in_cc := configuration{
		username:    "testuser",
		token:       "testtoken",
		ftEndpoint:  "testendpoint",
		ftClientDev: true,
		ftVersion:   "testversion",
	}
	//this will write a new .ft/creds.yaml file
	configToFile(testConfigFile, &in_cc)
	out_cc := configuration{}
	hydrateConfig(testConfigFile, &out_cc)
	assert.Equal(t, out_cc.ftEndpoint, "testendpoint")
	assert.Equal(t, out_cc.ftVersion, "testversion")
	assert.Equal(t, out_cc.ftClientDev, true)
	assert.Equal(t, out_cc.username, "testuser")
	assert.Equal(t, out_cc.token, "testtoken")
}

func TestConfigToMap(t *testing.T) {

	m := configToMap(&configuration{
		username:    "a",
		token:       "b",
		ftEndpoint:  "c",
		ftClientDev: true,
		ftVersion:   "0"})

	assert.Equal(t, m["ftEndpoint"], "c")
	assert.Equal(t, m["ftClientDev"], "true")
	assert.Equal(t, m["ftVersion"], "0")
	assert.Equal(t, m["username"], "a")
	assert.Equal(t, m["token"], "b")

}

func TestLoginCommand(t *testing.T) {
	cmd := exec.Command("../ft", "login", "-u", "usr", "-t", "tok")
	err := cmd.Start()
	cmd.Wait()
	assert.Nil(t, err)
	// creds should now be set:
	c := defaultConfig()
	assert.Equal(t, c.username, "usr")
	assert.Equal(t, c.token, "tok")

	//
	cmd = exec.Command("../ft", "login", "-u", "usr2", "-t", "tok2", "-e", "anEndpoint")
	err = cmd.Start()
	cmd.Wait()
	assert.Nil(t, err)
	// creds should now be set:
	c = defaultConfig()
	assert.Equal(t, c.username, "usr2")
	assert.Equal(t, c.token, "tok2")
	assert.Equal(t, c.ftEndpoint, "anEndpoint")

}
