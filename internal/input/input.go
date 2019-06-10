package input

import (
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type HandlerFunc func(*http.Request) (types.Validation, error)

type Input interface {
	Handler(*http.Request) (types.Validation, error)
}
