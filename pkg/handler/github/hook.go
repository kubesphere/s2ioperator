package github

import (
	"net/http"

	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
)

func Register(handlers []*builder.HandlerBuilder) {
	handlers = append(handlers, &builder.HandlerBuilder{
		Pattern: "/github",
		Func:    Serve,
	})
}

func Serve(w http.ResponseWriter, r *http.Request) {

}
