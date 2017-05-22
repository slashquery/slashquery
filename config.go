package slashquery

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
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
	s.Servers = make(map[string]Servers)
	return &s, nil
}
