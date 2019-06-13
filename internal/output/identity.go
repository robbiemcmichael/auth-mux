package output

import (
	"encoding/json"
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Identity struct{}

func (o *Identity) Handler(w http.ResponseWriter, validation types.Validation) error {
	return json.NewEncoder(w).Encode(validation)
}
