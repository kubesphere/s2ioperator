package handler

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
	"github.com/kubesphere/s2ioperator/pkg/handler/github"
)

var handlers = []*builder.HandlerBuilder{}

func init() {
	github.Register(handlers)
}
func Run() {
	for _, handler := range handlers {
		http.HandleFunc(handler.Pattern, handler.Func)
	}
	glog.Fatal(http.ListenAndServe(":8080", nil))
}
