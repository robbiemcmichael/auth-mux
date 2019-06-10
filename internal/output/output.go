package output

import (
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type HandlerFunc func(http.ResponseWriter, types.Validation) error

type Output interface {
	Handler(http.ResponseWriter, types.Validation) error
}
