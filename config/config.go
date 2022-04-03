package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func Parse(configFile string, config interface{}) error {
	content, err := ioutil.ReadFile(configFile)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, config)

	if err != nil {
		return err
	}

	return nil
}
