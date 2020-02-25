package github

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/google/go-github/github"
)

func (t Trigger) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"s2i_github_trigger"}
	ws.Path("/s2i.trigger/github")

	ws.Route(ws.POST("/namespaces/{namespace}/s2ibuilders/{s2ibuilder}").
		To(t.Serve).
		Consumes("application/x-www-form-urlencoded", "application/json", "charset=utf-8").
		Doc("trigger github handler").
		Param(ws.PathParameter("namespace", "namespace")).
		Param(ws.PathParameter("s2ibuilder", "the name of s2ibuilder")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(github.PushEvent{}))

	return ws
}
