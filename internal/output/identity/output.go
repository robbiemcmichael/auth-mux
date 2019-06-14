package identity

import (
	"encoding/json"
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Output struct{}

func (o *Output) Handler(w http.ResponseWriter, validation types.Validation) error {
	return json.NewEncoder(w).Encode(validation)
}
