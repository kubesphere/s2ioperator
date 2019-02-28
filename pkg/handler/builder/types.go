package builder

import (
	"net/http"
)

type HandlerBuilder struct {
	Pattern string
	Func    http.HandlerFunc
}
