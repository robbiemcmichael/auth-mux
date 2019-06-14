package token

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/robbiemcmichael/auth-mux/internal/token/jwt"
	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Validator interface {
	Validate(string) (types.Validation, error)
}

type Config struct {
	Type   string
	Config Validator
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var wrapper struct {
		Type   string      `yaml:"type"`
		Config interface{} `yaml:"config"`
	}

	if err := unmarshal(&wrapper); err != nil {
		return fmt.Errorf("unmarshal Token: %v", err)
	}

	var config Validator

	switch t := wrapper.Type; t {
	case "JWT":
		config = new(jwt.Config)
	default:
		return fmt.Errorf("unmarshal Token: unknown type %q", t)
	}

	b, err := yaml.Marshal(wrapper.Config)
	if err != nil {
		return fmt.Errorf("re-marshal config: %v", err)
	}

	// Unmarshal the token config based on the token type
	if err := yaml.Unmarshal(b, config); err != nil {
		return fmt.Errorf("unmarshal %s config: %v", wrapper.Type, err)
	}

	*c = Config{
		Type:   wrapper.Type,
		Config: config,
	}

	return nil
}
