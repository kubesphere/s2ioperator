package general

import (
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Test general webhook", func() {

	It("Should get s2irun after general webhook triggered", func() {
		res, err := http.Get(defaultUrl)
		Expect(err).NotTo(HaveOccurred(), "Cannot request handler with default url")
		Expect(res.StatusCode).To(Equal(http.StatusCreated))

		s2iruns := &devopsv1alpha1.S2iRunList{}

		err = t.KubeClientSet.List(context.TODO(), s2iruns, client.InNamespace(t.Namespace))
		Expect(err).NotTo(HaveOccurred(), "Can not get s2irun after general webhook triggered")

		instance := s2iruns.Items[0]
		Expect(instance.Spec.BuilderName).To(Equal(s2ibName))
	})
})
