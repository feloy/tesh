package mocks

import (
	"io"

	"go.yaml.in/yaml/v2"
)

type Mocks struct {
	Scenarios []Scenario
}

type Scenario struct {
	ID          string
	Description string
	Mocks       []Mock
}

type Mock struct {
	Description string
	Command     string
	Args        []string
	ExitCode    *int `yaml:"exit-code"`
	Stdout      *string
	Stderr      *string
}

func ParseMocks(mocksFile io.Reader) (*Mocks, error) {
	mocks := &Mocks{}
	err := yaml.NewDecoder(mocksFile).Decode(mocks)
	if err != nil {
		return nil, err
	}
	return mocks, nil
}
