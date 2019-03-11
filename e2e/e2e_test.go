package e2e_test

import (
	"context"
	"fmt"
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
	It("Should work well when using exactly the example yamls", func() {
		//create a s2ibuilder
		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err := os.Open(workspace + "/config/samples/devops_v1alpha1_s2ibuilder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), s2ibuilder)
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/devops_v1alpha1_s2irun.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), s2irun)

		var depKey = types.NamespacedName{Name: s2irun.Name + "-job", Namespace: s2irun.Namespace}
		var cmKey = types.NamespacedName{Name: s2irun.Name + "-configmap", Namespace: s2irun.Namespace}
		//configmap
		cm := &corev1.ConfigMap{}
		Eventually(func() error { return testClient.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		Expect(testClient.Delete(context.TODO(), cm)).NotTo(HaveOccurred())
		Eventually(func() error { return testClient.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		defer testClient.Delete(context.TODO(), cm)

		job := &batchv1.Job{}
		Eventually(func() error { return testClient.Get(context.TODO(), depKey, job) }, timeout, time.Second).
			Should(Succeed())

		//for our example, the status must be successful
		Eventually(func() error {
			err = testClient.Get(context.TODO(), depKey, job)
			if err != nil {
				return err
			}
			if job.Status.Succeeded == 1 {
				//Status of s2ibuilder should update too
				tempBuilder := &devopsv1alpha1.S2iBuilder{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, tempBuilder)
				if err != nil {
					return err
				}
				if *tempBuilder.Status.LastRunName == s2irun.Name && tempBuilder.Status.LastRunState == devopsv1alpha1.Successful && tempBuilder.Status.RunCount == 1 {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())

		// Delete the Deployment and expect Reconcile to be called for Deployment deletion
		Expect(testClient.Delete(context.TODO(), job)).NotTo(HaveOccurred())
		Eventually(func() error { return testClient.Get(context.TODO(), depKey, job) }, timeout).
			Should(Succeed())
	})
})
