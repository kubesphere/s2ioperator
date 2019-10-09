package e2e_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = Describe("", func() {

	const (
		timeout           = time.Second * 25
		TaintKey          = "node.kubernetes.io/ci"
		NodeAffinityKey   = "node-role.kubernetes.io/worker"
		NodeAffinityValue = "ci"
	)
	It("Should work well when using exactly the runtimeimage example yamls", func() {
		//create a s2ibuilder
		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err := os.Open(workspace + "/config/samples/devops_v1alpha2_s2ibuilder_runtimeimage.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/devops_v1alpha2_s2irun_runtimeimage.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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

		job := &batchv1.Job{}
		Eventually(func() error { return testClient.Get(context.TODO(), depKey, job) }, timeout, time.Second).
			Should(Succeed())

		res := checkAnffinitTaint(job, NodeAffinityKey, NodeAffinityValue, TaintKey)
		Expect(res).To(Equal(true))

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
		}, time.Minute*10, time.Second*10).Should(Succeed())

		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())
	})

	It("Should work well when using exactly the example yamls", func() {
		//create a s2ibuilder
		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err := os.Open(workspace + "/config/samples/devops_v1alpha1_s2ibuilder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/devops_v1alpha1_s2irun.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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

		job := &batchv1.Job{}
		Eventually(func() error { return testClient.Get(context.TODO(), depKey, job) }, timeout, time.Second).
			Should(Succeed())

		res := checkAnffinitTaint(job, s2ibuilder.Spec.Config.NodeAffinityKey, s2ibuilder.Spec.Config.NodeAffinityValues[0], s2ibuilder.Spec.Config.TaintKey)
		Expect(res).To(Equal(true))

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

		Eventually(func() bool {
			res := &devopsv1alpha1.S2iRun{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, res)
			if err != nil {
				return false
			}
			if strings.Contains(res.Status.S2iBuildResult.ImageName, s2ibuilder.Spec.Config.ImageName) {
				return true
			}
			return false
		}, timeout, time.Second).Should(BeTrue())

		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())
	})

	It("Should autoScale work well when using exactly the example yamls", func() {
		//create a s2ibuilder
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
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/autoscale/python-s2i-run.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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
			return fmt.Errorf("Failed %+v", job)
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
		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {

			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())

		Eventually(func() error {

			testStatefulSet := &appsv1.StatefulSet{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, testStatefulSet)
			if err != nil {
				return err
			}
			if _, ok := testStatefulSet.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove statefulset label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())

		Eventually(func() error {

			testDeployment := &appsv1.Deployment{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, testDeployment)
			if err != nil {
				return err
			}
			if _, ok := testDeployment.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove statefulset label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), statefulSet) }, timeout, time.Second).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), deploy) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, deploy))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, statefulSet))
		}, timeout, time.Second).Should(BeTrue())
	})

	It("Should autoScale fail when using exactly the example yamls", func() {
		//create a s2ibuilder
		deploy := &appsv1.Deployment{}

		reader, err := os.Open(workspace + "/config/samples/autoscale-failed/python-deployment.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(deploy)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), deploy)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), deploy)

		statefulSet := &appsv1.StatefulSet{}
		reader, err = os.Open(workspace + "/config/samples/autoscale-failed/python-statefulset.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(statefulSet)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), statefulSet)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), statefulSet)

		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err = os.Open(workspace + "/config/samples/autoscale-failed/python-s2i-builder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/autoscale-failed/python-s2i-run.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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
		job := &batchv1.Job{}
		Eventually(func() error { return testClient.Get(context.TODO(), depKey, job) }, timeout, time.Second).
			Should(Succeed())

		//for our example, the status must be successful
		Eventually(func() error {
			err = testClient.Get(context.TODO(), depKey, job)
			if err != nil {
				return err
			}
			if job.Status.Failed == 1 {
				//Status of s2ibuilder should update too
				tempBuilder := &devopsv1alpha1.S2iBuilder{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, tempBuilder)
				if err != nil {
					return err
				}
				if tempBuilder.Status.LastRunName != nil && *tempBuilder.Status.LastRunName == s2irun.Name && tempBuilder.Status.LastRunState == devopsv1alpha1.Failed && tempBuilder.Status.RunCount == 1 {
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
			if job.Status.Failed == 1 {
				//Status of s2ibuilder should update too
				testDeploy := &appsv1.Deployment{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, testDeploy)
				if err != nil {
					return err
				}
				if testDeploy.Annotations[devopsv1alpha1.WorkLoadCompletedInitAnnotations] == devopsv1alpha1.Failed {
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
			if job.Status.Failed == 1 {
				//Status of s2ibuilder should update too
				testStatefulSet := &appsv1.StatefulSet{}
				err = testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, testStatefulSet)
				if err != nil {
					return err
				}
				if testStatefulSet.Annotations[devopsv1alpha1.WorkLoadCompletedInitAnnotations] == devopsv1alpha1.Failed {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())

		Eventually(func() error {

			testStatefulSet := &appsv1.StatefulSet{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, testStatefulSet)
			if err != nil {
				return err
			}
			if _, ok := testStatefulSet.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove statefulset label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())

		Eventually(func() error {

			testDeployment := &appsv1.Deployment{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, testDeployment)
			if err != nil {
				return err
			}
			if _, ok := testDeployment.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove deploy label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), statefulSet) }, timeout, time.Second).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), deploy) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, deploy))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, statefulSet))
		}, timeout, time.Second).Should(BeTrue())
	})

	It("Should work well when using secrets", func() {
		//create a s2ibuilder
		secret := &corev1.Secret{}
		reader, err := os.Open(workspace + "/config/samples/secret/secret.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(secret)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), secret)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), secret)

		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err = os.Open(workspace + "/config/samples/secret/devops_v1alpha1_s2ibuilder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/secret/devops_v1alpha1_s2irun.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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
		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())

	})

	It("Should work well when using git secrets", func() {
		//create a s2ibuilder
		secret := &corev1.Secret{}
		reader, err := os.Open(workspace + "/config/samples/git-secret/secret.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(secret)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), secret)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), secret)

		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err = os.Open(workspace + "/config/samples/git-secret/devops_v1alpha1_s2ibuilder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/git-secret/devops_v1alpha1_s2irun.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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
		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())
	})

	It("Should b2i work well when using exactly the example yamls", func() {
		//create a s2ibuilder
		deploy := &appsv1.Deployment{}

		reader, err := os.Open(workspace + "/config/samples/b2i/tomcat-deployment.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(deploy)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), deploy)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), deploy)

		statefulSet := &appsv1.StatefulSet{}
		reader, err = os.Open(workspace + "/config/samples/b2i/tomcat-statefulset.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(statefulSet)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), statefulSet)
		Expect(err).NotTo(HaveOccurred())
		defer testClient.Delete(context.TODO(), statefulSet)

		s2ibuilder := &devopsv1alpha1.S2iBuilder{}
		reader, err = os.Open(workspace + "/config/samples/b2i/tomcat-s2i-builder.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2ibuilder)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		s2irun := &devopsv1alpha1.S2iRun{}
		reader, err = os.Open(workspace + "/config/samples/b2i/tomcat-s2i-run.yaml")
		Expect(err).NotTo(HaveOccurred(), "Cannot read sample yamls")
		err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(s2irun)
		Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yamls")
		err = testClient.Create(context.TODO(), s2irun)
		Expect(err).NotTo(HaveOccurred())

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
				if testDeploy.Spec.Replicas != nil && *testDeploy.Spec.Replicas == 3 && testDeploy.Spec.Template.Spec.Containers[0].Image == "runzexia/hello-java:v0.1.0" {
					return nil
				}
			}
			return fmt.Errorf("Failed %+v", job)
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
				if testStatefulSet.Spec.Replicas != nil && *testStatefulSet.Spec.Replicas == 1 && testStatefulSet.Spec.Template.Spec.Containers[0].Image == "runzexia/hello-java:v0.1.0" {
					return nil
				}
			}
			return fmt.Errorf("Failed")
		}, time.Minute*5, time.Second*10).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), s2ibuilder) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {

			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2irun.Name, Namespace: s2irun.Namespace}, s2irun))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: s2ibuilder.Name, Namespace: s2ibuilder.Namespace}, s2ibuilder))
		}, timeout, time.Second).Should(BeTrue())

		Eventually(func() error {

			testStatefulSet := &appsv1.StatefulSet{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, testStatefulSet)
			if err != nil {
				return err
			}
			if _, ok := testStatefulSet.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove statefulset label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())

		Eventually(func() error {

			testDeployment := &appsv1.Deployment{}
			err = testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, testDeployment)
			if err != nil {
				return err
			}
			if _, ok := testDeployment.Labels[s2ibuilder.Name]; ok {
				return fmt.Errorf("should remove statefulset label")
			}
			return nil
		}, time.Minute*5, time.Second*10).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), statefulSet) }, timeout, time.Second).Should(Succeed())
		Eventually(func() error { return testClient.Delete(context.TODO(), deploy) }, timeout, time.Second).Should(Succeed())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, deploy))
		}, timeout, time.Second).Should(BeTrue())
		Eventually(func() bool {
			return errors.IsNotFound(testClient.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, statefulSet))
		}, timeout, time.Second).Should(BeTrue())
	})
})

func checkAnffinitTaint(job *batchv1.Job, nodeAffinityKey string, nodeAffinityValue string, taintKey string) bool {
	if job.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].Preference.MatchExpressions[0].Key != nodeAffinityKey {
		fmt.Println("Check nodeAffinityKey error.")
		fmt.Println(job.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].Preference.MatchExpressions[0].Key, nodeAffinityKey)
		return false

	}
	if job.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].Preference.MatchExpressions[0].Values[0] != nodeAffinityValue {
		fmt.Println("Check nodeAffinityValue error.")
		fmt.Println(job.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].Preference.MatchExpressions[0].Values[0], nodeAffinityValue)
		return false
	}

	for _, toleration := range job.Spec.Template.Spec.Tolerations {
		if toleration.Key == taintKey {
			return true
		}
	}
	fmt.Println("Check toleration error.")
	return false
}
