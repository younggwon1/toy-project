package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func readFromFile(filename string, out interface{}) error {
	// open file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// unmarshal data into cfg
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
