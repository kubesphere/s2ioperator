package e2e_test

import (
	"context"
	"os"
	"time"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = Describe("", func() {

	const timeout = time.Second * 25
	It("Should work well", func() {
		//create a s2ibuilder
		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err := os.Open(workspace + "/config/samples/devops_v1alpha1_s2ibuilder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = c.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), s2ibuilder)
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/devops_v1alpha1_s2irun.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = c.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), s2irun)

		var depKey = types.NamespacedName{Name: s2irun.Name + "-job", Namespace: s2irun.Namespace}
		var cmKey = types.NamespacedName{Name: s2irun.Name + "-configmap", Namespace: s2irun.Namespace}
		//configmap
		cm := &corev1.ConfigMap{}
		Eventually(func() error { return c.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		Expect(c.Delete(context.TODO(), cm)).NotTo(HaveOccurred())
		Eventually(func() error { return c.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		defer c.Delete(context.TODO(), cm)

		job := &batchv1.Job{}
		Eventually(func() error { return c.Get(context.TODO(), depKey, job) }, timeout).
			Should(Succeed())

		// Delete the Deployment and expect Reconcile to be called for Deployment deletion
		Expect(c.Delete(context.TODO(), job)).NotTo(HaveOccurred())
		Eventually(func() error { return c.Get(context.TODO(), depKey, job) }, timeout).
			Should(Succeed())
	})
})
