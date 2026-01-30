package scenarios

import (
	"io"

	"go.yaml.in/yaml/v2"
)

type Scenarios struct {
	Scenarios []Scenario
}

type Scenario struct {
	ID          string
	Description string
	Mocks       []Mock
	Envs        []string
	Expect      *Expect
}

type Mock struct {
	Description string
	Command     string
	Args        []string
	ExitCode    *int `yaml:"exit-code"`
	Stdout      *string
	Stderr      *string
}

type Expect struct {
	ExitCode *int `yaml:"exit-code"`
	Stdout   *string
	Stderr   *string
	Calls    []Call
}

type Call struct {
	Command string
	Args    []string
	Called  int
}

func ParseScenarios(scenariosFile io.Reader) (*Scenarios, error) {
	scenarios := &Scenarios{}
	err := yaml.NewDecoder(scenariosFile).Decode(scenarios)
	if err != nil {
		return nil, err
	}
	return scenarios, nil
}
