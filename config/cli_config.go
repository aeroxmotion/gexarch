package config

import (
	"io/ioutil"

	"github.com/aeroxmotion/gexarch/util"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILENAME = "gexarch.yml"

type CliConfig struct {
	TargetPath string `yaml:"target_path"`
}

func GetCliConfig() *CliConfig {
	cliConfig := CliConfig{}

	configBytes, err := ioutil.ReadFile(CONFIG_FILENAME)
	util.PanicIfError(err)

	util.PanicIfError(yaml.Unmarshal(configBytes, &cliConfig))

	return &cliConfig
}
