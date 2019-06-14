package kubernetesTokenReview

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "k8s.io/api/authentication/v1"

	"github.com/robbiemcmichael/auth-mux/internal/token"
	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Config struct {
	Validator token.Config `yaml:"validator"`
}

func (c *Config) Handler(r *http.Request) (types.Validation, error) {
	decoder := json.NewDecoder(r.Body)

	var tokenReview auth.TokenReview
	if err := decoder.Decode(&tokenReview); err != nil {
		return types.Validation{}, fmt.Errorf("decode JSON: %v", err)
	}

	validation, err := c.Validator.Config.Validate(tokenReview.Spec.Token)
	if err != nil {
		return types.Validation{}, fmt.Errorf("validate token: %v", err)
	}

	return validation, nil
}
