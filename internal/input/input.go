package input

import (
	"net/http"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type HandlerFunc func(*http.Request) (types.Result, error)

type Input interface {
	Handler(*http.Request) (types.Result, error)
}
