package general

import (
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/client/clientset/versioned/scheme"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	t Trigger
)

const (
	defaultUrl = "http://127.0.0.1:8000/s2itrigger/v1alpha1/general/namespaces/" + namespace + "/s2ibuilders/" + s2ibName
	s2ibName   = "s2i-b"
	namespace  = "s2i"
)

func TestS2irun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "General webhook Suite")
}

var _ = BeforeSuite(func() {
	s2ib := &devopsv1alpha1.S2iBuilder{
		ObjectMeta: v1.ObjectMeta{
			Name:      s2ibName,
			Namespace: namespace,
		},
		Spec: devopsv1alpha1.S2iBuilderSpec{
			Config: &devopsv1alpha1.S2iConfig{
				RevisionId: "master",
				SecretCode: "secretCode",
			},
		},
	}

	scheme := scheme.Scheme
	c := fake.NewFakeClientWithScheme(scheme, s2ib)
	t.KubeClientSet = c
	t.Namespace = s2ib.Namespace
	t.S2iBuilderName = s2ib.Name
})
