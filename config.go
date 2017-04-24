package slashquery

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// New return a new slashquery instance
func New(file string) (*Slashquery, error) {
	ymlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var s Slashquery
	if err := yaml.Unmarshal(ymlFile, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
