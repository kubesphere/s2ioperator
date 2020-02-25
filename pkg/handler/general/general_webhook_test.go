package general

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test general webhook", func() {

	It("Should get s2irun after general webhook triggered with get request", func() {
		ws := new(restful.WebService)
		ws.Path("/s2itrigger/v1alpha1/general")
		ws.Route(ws.GET("/namespaces/{namespace}/s2ibuilders/{s2ibuilder}").
			To(t.Serve).
			Doc("trigger general handler with GET").
			Param(ws.PathParameter("namespace", "namespace")).
			Param(ws.PathParameter("s2ibuilder", "the name of s2ibuilder")).
			Param(ws.QueryParameter("secretCode", "use secret code to authorizing").
				Required(true).
				DataFormat("secretCode=%s")).
			Metadata(restfulspec.KeyOpenAPITags, tags))
		restful.Add(ws)

		reqUrl := defaultUrl + "?secretCode=secretCode"
		httpRequest, _ := http.NewRequest("GET", reqUrl, nil)
		httpWriter := httptest.NewRecorder()

		restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)
		Expect(httpWriter.Code).To(Equal(http.StatusCreated))

		s2iruns := &devopsv1alpha1.S2iRunList{}

		err := t.KubeClientSet.List(context.TODO(), s2iruns, client.InNamespace(t.Namespace))
		Expect(err).NotTo(HaveOccurred(), "Can not get s2irun after general webhook triggered")

		instance := s2iruns.Items[0]
		Expect(instance.Spec.BuilderName).To(Equal(s2ibName))

		t.KubeClientSet.Delete(context.TODO(), &instance)
	})

	It("Should trigger failed without secretCode", func() {

		httpRequest, _ := http.NewRequest("GET", defaultUrl, nil)
		httpWriter := httptest.NewRecorder()
		restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

		Expect(httpWriter.Code).To(Equal(http.StatusUnauthorized))
	})
})
