package general

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
)

var tags = []string{"s2i_general_trigger"}

func (t Trigger) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/s2itrigger/v1alpha1/general")

	// handle request with GET method
	ws.Route(ws.GET("/namespaces/{namespace}/s2ibuilders/{s2ibuilder}").
		To(t.Serve).
		Doc("trigger general handler with GET").
		Param(ws.PathParameter("namespace", "namespace")).
		Param(ws.PathParameter("s2ibuilder", "the name of s2ibuilder")).
		Param(ws.QueryParameter("secretCode", "use secret code to authorizing").
			Required(true).
			DataFormat("secretCode=%s")).
		Metadata(restfulspec.KeyOpenAPITags, tags))

	// handle request with GET method
	ws.Route(ws.POST("/namespaces/{namespace}/s2ibuilders/{s2ibuilder}").
		To(t.Serve).
		Doc("trigger general handler with POST").
		Consumes("application/x-www-form-urlencoded", "application/json", "charset=utf-8").
		Param(ws.PathParameter("namespace", "namespace")).
		Param(ws.PathParameter("s2ibuilder", "the name of s2ibuilder")).
		Param(ws.QueryParameter("secretCode", "use secret code to authorizing").
			Required(true).
			DataFormat("secretCode=%s")).
		Metadata(restfulspec.KeyOpenAPITags, tags))

	return ws
}
