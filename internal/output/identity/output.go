package identity

import (
	"encoding/json"
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Config struct{}

func (c *Config) Handler(w http.ResponseWriter, validation types.Validation) error {
	if !validation.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}

	return json.NewEncoder(w).Encode(validation)
}
