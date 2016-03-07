package flowthings

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/gopkg.in/yaml.v2"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/flowthings/ft-plan-cli/utils"
)

// Configuration is a struct that stores a user's configuraiton
type configuration struct {
	username    string
	token       string
	ftEndpoint  string
	ftClientDev bool
	ftVersion   string
	configFile  string
	planFile    string
	basePath    string
	planDir     string
	tempDir     string
	planID      string
}

// errIfInvalid looks at whether or not the username and password is valid.
// And it will error out if there's not one or the other in the config.
func errIfInvalid(c *configuration) error {
	if c.username == "" && c.token == "" {
		return errors.New("No username is configured. No token is configured.")
	}
	if c.username == "" {
		return errors.New("No username is configured")
	}
	if c.token == "" {
		return errors.New("No token is configured")
	}
	return nil
}

func updateConfig(from *configuration, to *configuration) {
	//hydrate config from specified file if it exists or else default
	//if it exists. After that copy over any other manually specified
	//values
	configFile := from.configFile
	if configFile == "" {
		configFile = defaultConfigFile()
	}
	if !fileExists(configFile) {
		utils.Checkerror(fmt.Errorf("can't find config file %s", configFile))
	}
	hydrateGlobalConfig(configFile, to)
	//copy over manually specified values
	localConfigFile := filepath.Join(from.planDir, "plan.yml")
	if fileExists(localConfigFile) {
		to.planFile = localConfigFile
		hydrateLocalConfig(localConfigFile, to)
	}

	if from.username != "" {
		to.username = from.username
	}
	if from.token != "" {
		to.token = from.token
	}
	if from.ftEndpoint != "" && from.ftEndpoint != defaultPlatformEndpoint {
		to.ftEndpoint = from.ftEndpoint
	}
	if from.ftClientDev {
		to.ftClientDev = true
	}
	if from.basePath != "" {
		to.basePath = from.basePath
	}
	if from.planDir != "" {
		to.planDir = from.planDir
	}
	to.tempDir = createTempDirectory(path.Dir(configFile), to.ftClientDev)
}

// hydrateGlobalConfig loads the general config,
// then it loads the username and token into the configuration variable.
//fill in config value with file contents if they exist.
//caller must ensure that file exists
func hydrateGlobalConfig(fullFileName string, c *configuration) {
	v := hydrateConfig(fullFileName, c)

	c.configFile = fullFileName
	c.username = v.GetString("username")
	c.token = v.GetString("token")
	c.ftClientDev = v.GetBool("ft_client_dev")
}

// hydrateLocalConfig loads the general config and then loads the basePath
func hydrateLocalConfig(fullFileName string, c *configuration) {
	v := hydrateConfig(fullFileName, c)
	c.basePath = v.GetString("base_path")
	c.planID = v.GetString("plan_id")
}

// hydrateConfig is the general configuration shared by local and global configuration
func hydrateConfig(fullFileName string, c *configuration) *viper.Viper {
	v := viper.New()
	//set defaults
	v.SetDefault("ft_endpoint", defaultPlatformEndpoint)
	v.SetDefault("ft_version", defaultPlatformVersion)
	v.SetDefault("ft_client_dev", false)

	// set file
	fileDirName := path.Dir(fullFileName)
	filename := path.Base(fullFileName)
	fileBaseName := strings.Split(filename, ".")[0] //file.ext -> file
	v.SetConfigName(fileBaseName)
	v.AddConfigPath(fileDirName)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	utils.Checkerror(err)

	c.ftEndpoint = v.GetString("ft_endpoint")
	c.ftVersion = v.GetString("ft_version")

	return v
}

// configToFile looks at the config and then writes it to the file.
func configToFile(filename string, c *configuration) {
	configMap := configToMap(c)

	if configMap["ft_endpoint"] == defaultPlatformEndpoint || configMap["ft_endpoint"] == "" {
		delete(configMap, "ft_endpoint")
	}
	if configMap["ft_client_dev"] == "false" || configMap["ft_client_dev"] == "" {
		delete(configMap, "ft_client_dev")
	}
	if configMap["ft_version"] == defaultPlatformVersion || configMap["ft_version"] == "" {
		delete(configMap, "ft_version")
	}

	configDir := path.Dir(filename)
	if err := os.Mkdir(configDir, 0777); err != nil && !os.IsExist(err) {
		utils.Checkerror(fmt.Errorf("%s: \"%s\"", err, configDir))
	}

	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		utils.Checkerror(err)
	}

	file, err := os.Create(filename)
	if err != nil && !os.IsNotExist(err) {
		utils.Checkerror(err)
	}

	yaml, err := yaml.Marshal(&configMap)
	if _, err := file.Write(yaml); err != nil {
		utils.Checkerror(fmt.Errorf("%s: %s", err, filename))
	}
}

func configToMap(c *configuration) map[string]string {
	yconfig := make(map[string]string)
	yconfig["ft_endpoint"] = c.ftEndpoint
	yconfig["ft_client_dev"] = strconv.FormatBool(c.ftClientDev)
	yconfig["ft_version"] = c.ftVersion
	yconfig["username"] = c.username
	yconfig["token"] = c.token
	return yconfig
}

// get config Directory @ $HOME/.ft. Error if the directory is not accessible or writable
func homeConfigDirectory() string {
	home := os.Getenv("HOME")

	if home == "" {
		utils.Checkerror(errors.New("Can't read HOME environment variable."))
	}

	configDir := path.Join(home, ".ft")

	if err := os.Mkdir(configDir, 0777); err != nil && !os.IsExist(err) {
		utils.Checkerror(fmt.Errorf("Can't access config dir %s", configDir))
	}

	return configDir
}

func defaultConfig() *configuration {
	cwd, err := os.Getwd()
	utils.Checkerror(err)
	c := configuration{
		planDir:     cwd,
		ftEndpoint:  defaultPlatformEndpoint,
		ftVersion:   defaultPlatformVersion,
		ftClientDev: false,
	}

	configFile := defaultConfigFile()
	if fileExists(configFile) {
		hydrateGlobalConfig(configFile, &c)
	}
	return &c
}

// ------------------- supporting -----------------------

func planURI(c *configuration) string {
	return strings.Join([]string{c.ftEndpoint, c.ftVersion, c.username, "plan"}, "/")
}

func defaultConfigFile() string {
	return path.Join(homeConfigDirectory(), "creds.yaml")
}
