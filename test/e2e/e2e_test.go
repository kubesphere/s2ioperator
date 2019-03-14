package e2e_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = Describe("", func() {

	const timeout = time.Second * 25
	It("Should work well when using exactly the example yamls", func() {
		//create a s2ibuilder
		cleanDelete := client.PropagationPolicy(metav1.DeletePropagationBackground)
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
		defer testClient.Delete(context.TODO(), s2irun, cleanDelete)

		createdInstance := &devopsv1alpha1.S2iRun{}
		Eventually(func() error {
			return testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, createdInstance)
		}, timeout).Should(Succeed())

		instanceUidSlice := strings.Split(string(createdInstance.UID), "-")
		var cmKey = types.NamespacedName{
			Name:      s2irun.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-configmap",
			Namespace: s2irun.Namespace}
		var depKey = types.NamespacedName{
			Name:      s2irun.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-job",
			Namespace: s2irun.Namespace}
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
				if tempBuilder.Status.LastRunName != nil && *tempBuilder.Status.LastRunName == s2irun.Name && tempBuilder.Status.LastRunState == devopsv1alpha1.Successful && tempBuilder.Status.RunCount == 1 {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())
	})

	It("Should autoScale work well well when using exactly the example yamls", func() {
		//create a s2ibuilder
		cleanDelete := client.PropagationPolicy(metav1.DeletePropagationBackground)
		deploy := &appsv1.Deployment{}

		reader, err := os.Open(workspace + "/config/samples/autoscale/python-deployment.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(deploy)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), deploy)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), deploy)

		statefulSet := &appsv1.StatefulSet{}
		reader, err = os.Open(workspace + "/config/samples/autoscale/python-statefulset.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(statefulSet)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), statefulSet)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), statefulSet)

		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err = os.Open(workspace + "/config/samples/autoscale/python-s2i-builder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), s2ibuilder)
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/autoscale/python-s2i-run.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), s2irun, cleanDelete)

		createdInstance := &devopsv1alpha1.S2iRun{}
		Eventually(func() error {
			return testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, createdInstance)
		}, timeout).Should(Succeed())

		instanceUidSlice := strings.Split(string(createdInstance.UID), "-")
		var cmKey = types.NamespacedName{
			Name:      s2irun.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-configmap",
			Namespace: s2irun.Namespace}
		var depKey = types.NamespacedName{
			Name:      s2irun.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-job",
			Namespace: s2irun.Namespace}
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
				if tempBuilder.Status.LastRunName != nil && *tempBuilder.Status.LastRunName == s2irun.Name && tempBuilder.Status.LastRunState == devopsv1alpha1.Successful && tempBuilder.Status.RunCount == 1 {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())

		Eventually(func() error {
			err = testClient.Get(context.TODO(), depKey, job)
			if err != nil {
				return err
			}
			if job.Status.Succeeded == 1 {
				//Status of s2ibuilder should update too
				testDeploy := &appsv1.Deployment{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, testDeploy)
				if err != nil {
					return err
				}
				if testDeploy.Spec.Replicas != nil && *testDeploy.Spec.Replicas == 3 && testDeploy.Spec.Template.Spec.Containers[0].Image == "runzexia/hello-python:v0.1.0" {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())

		Eventually(func() error {
			err = testClient.Get(context.TODO(), depKey, job)
			if err != nil {
				return err
			}
			if job.Status.Succeeded == 1 {
				//Status of s2ibuilder should update too
				testStatefulSet := &appsv1.StatefulSet{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, testStatefulSet)
				if err != nil {
					return err
				}
				if testStatefulSet.Spec.Replicas != nil && *testStatefulSet.Spec.Replicas == 1 && testStatefulSet.Spec.Template.Spec.Containers[0].Image == "runzexia/hello-python:v0.1.0" {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())
	})

})
