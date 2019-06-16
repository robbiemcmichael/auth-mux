package bearer

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/robbiemcmichael/auth-mux/internal/token"
	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Config struct {
	Validator token.Config `yaml:"validator"`
}

func (c *Config) Handler(r *http.Request) (types.Validation, error) {
	header := r.Header.Get("authorization")

	if header == "" {
		invalid := types.Validation{
			Valid: false,
			Error: "Missing authorization header",
		}
		return invalid, nil
	}

	auth := strings.SplitN(header, " ", 2)

	if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
		invalid := types.Validation{
			Valid: false,
			Error: "Expected bearer token in authorization header",
		}
		return invalid, nil
	}

	validation, err := c.Validator.Config.Validate(auth[1])
	if err != nil {
		return types.Validation{}, fmt.Errorf("validate token: %v", err)
	}

	return validation, nil
}
