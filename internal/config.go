package internal

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/robbiemcmichael/auth-mux/internal/input"
	"github.com/robbiemcmichael/auth-mux/internal/output"
)

type Config struct {
	Address string   `yaml:"address"`
	Port    int      `yaml:"port"`
	Cert    string   `yaml:"cert"`
	Key     string   `yaml:"key"`
	Inputs  []Input  `yaml:"inputs"`
	Outputs []Output `yaml:"outputs"`
}

// An Input takes an HTTP request and returns a Validation
type Input struct {
	Type   string
	Name   string
	Path   string
	Config input.Input
}

func (i *Input) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var wrapper struct {
		Type   string      `yaml:"type"`
		Name   string      `yaml:"name"`
		Path   string      `yaml:"path"`
		Config interface{} `yaml:"config"`
	}

	if err := unmarshal(&wrapper); err != nil {
		return fmt.Errorf("unmarshal Input: %v", err)
	}

	var config input.Input

	switch t := wrapper.Type; t {
	case "KubernetesTokenReview":
		config = new(input.KubernetesTokenReview)
	default:
		return fmt.Errorf("unmarshal Input: unknown type %q", t)
	}

	b, err := yaml.Marshal(wrapper.Config)
	if err != nil {
		return fmt.Errorf("re-marshal config: %v", err)
	}

	// Unmarshal the input config based on the input type
	if err := yaml.Unmarshal(b, config); err != nil {
		return fmt.Errorf("unmarshal Input config: %v", err)
	}

	*i = Input{
		Type:   wrapper.Type,
		Name:   wrapper.Name,
		Path:   wrapper.Path,
		Config: config,
	}

	return nil
}

// An Ouput takes a Validation and returns an HTTP response
type Output struct {
	Type   string
	Name   string
	Path   string
	Config output.Output
}

func (o *Output) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var wrapper struct {
		Type   string      `yaml:"type"`
		Name   string      `yaml:"name"`
		Path   string      `yaml:"path"`
		Config interface{} `yaml:"config"`
	}

	if err := unmarshal(&wrapper); err != nil {
		return fmt.Errorf("unmarshal Output: %v", err)
	}

	var config output.Output

	switch t := wrapper.Type; t {
	case "Identity":
		config = new(output.Identity)
	case "KubernetesTokenReview":
		config = new(output.KubernetesTokenReview)
	default:
		return fmt.Errorf("unmarshal Output: unknown type %q", t)
	}

	b, err := yaml.Marshal(wrapper.Config)
	if err != nil {
		return fmt.Errorf("re-marshal config: %v", err)
	}

	// Unmarshal the output config based on the output type
	if err := yaml.Unmarshal(b, config); err != nil {
		return fmt.Errorf("unmarshal Output config: %v", err)
	}

	*o = Output{
		Type:   wrapper.Type,
		Name:   wrapper.Name,
		Path:   wrapper.Path,
		Config: config,
	}

	return nil
}
