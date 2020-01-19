package handler

import (
	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
	"github.com/kubesphere/s2ioperator/pkg/handler/gitlab"
)

func init() {
	handlers = append(handlers, &builder.HandlerBuilder{
		Pattern: "/gitlab/",
		Func:    gitlab.NewGitlabSink(builder.ClientSets()).Serve,
	})
}
