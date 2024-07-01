package service

import (
	"os"

	"gopkg.in/yaml.v3"
)

const filename = "svc.yml"

type Svc struct {
	Ports map[string]string `yaml:"ports"`
}

func ParseFile() (*Svc, error) {
	svcFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var svc Svc
	if err := yaml.Unmarshal(svcFile, &svc); err != nil {
		return nil, err
	}

	return &svc, nil
}
