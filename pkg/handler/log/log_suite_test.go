package log_test

import (
	"testing"

	"github.com/kubesphere/s2ioperator/pkg/handler/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Testing logging URL for kubesphere", func() {
	It("Should get right url", func() {
		logger := log.GetKubesphereLogger()
		str, err := logger.GetURLOfPodLog("default", "pod-1")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(str).To(Equal("http://ks-apigateway.kubesphere-system.svc/apis/logging.kubesphere.io/v1alpha2/namespaces/default/pods/pod-1?operation=query"))
	})
})
