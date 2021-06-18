/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package s2irun

import (
	"context"
	"fmt"
	"strings"
	"time"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Test reconcile", func() {

	var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}

	const timeout = time.Second * 10

	It("Should get job and configmap when everything is right", func() {
		instance := &devopsv1alpha1.S2iRun{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
			Spec: devopsv1alpha1.S2iRunSpec{
				BuilderName: "foo",
			},
		}
		//create a s2ibuilder

		s2ibuilder := &devopsv1alpha1.S2iBuilder{
			ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
			Spec: devopsv1alpha1.S2iBuilderSpec{
				Config: &devopsv1alpha1.S2iConfig{
					ImageName: "hello/world",
					Tag:       "latest",
				},
			},
		}
		err := c.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), s2ibuilder)
		// Create the S2iRun object and expect the Reconcile and Deployment to be created
		err = c.Create(context.TODO(), instance)
		// The instance object may not be a valid object because it might be missing some required fields.
		// Please modify the instance object by adding required fields and then remove the following if statement.
		if apierrors.IsInvalid(err) {
			fmt.Printf("failed to create object, got an invalid object error: %v", err)
			return
		}
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), instance)
		Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))

		createdInstance := &devopsv1alpha1.S2iRun{}
		Eventually(func() error {
			return c.Get(context.TODO(), types.NamespacedName{Name: "foo", Namespace: "default"}, createdInstance)
		}, timeout).Should(Succeed())
		instanceUidSlice := strings.Split(string(createdInstance.UID), "-")

		var depKey = types.NamespacedName{Name: instance.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-job", Namespace: "default"}
		var cmKey = types.NamespacedName{Name: instance.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-configmap", Namespace: "default"}
		//configmap
		cm := &corev1.ConfigMap{}
		Eventually(func() error { return c.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		Expect(c.Delete(context.TODO(), cm)).NotTo(HaveOccurred())
		Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
		Eventually(func() error { return c.Get(context.TODO(), cmKey, cm) }, timeout).
			Should(Succeed())
		defer c.Delete(context.TODO(), cm)

		job := &batchv1.Job{}
		Eventually(func() error { return c.Get(context.TODO(), depKey, job) }, timeout).
			Should(Succeed())

		// Delete the Deployment and expect Reconcile to be called for Deployment deletion
		Expect(c.Delete(context.TODO(), job)).NotTo(HaveOccurred())
		Eventually(requests, timeout).Should(Receive(Equal(expectedRequest)))
		Eventually(func() error { return c.Get(context.TODO(), depKey, job) }, timeout).
			Should(Succeed())

		// Manually delete Deployment since GC isn't enabled in the test control plane
		Expect(c.Delete(context.TODO(), job)).To(Succeed())
	})
	It("Should not config from nonexist templates", func() {
		instance := &devopsv1alpha1.S2iRun{ObjectMeta: metav1.ObjectMeta{Name: "foo1", Namespace: "default"},
			Spec: devopsv1alpha1.S2iRunSpec{
				BuilderName: "foo1",
			},
		}
		//create a s2ibuildern using template
		s2ibuilder := &devopsv1alpha1.S2iBuilder{
			ObjectMeta: metav1.ObjectMeta{Name: "foo1", Namespace: "default"},
			Spec: devopsv1alpha1.S2iBuilderSpec{
				FromTemplate: &devopsv1alpha1.UserDefineTemplate{
					Name: "NotExsit",
				},
				Config: &devopsv1alpha1.S2iConfig{
					Tag: "latest",
				},
			},
		}
		err := c.Create(context.TODO(), s2ibuilder)
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), s2ibuilder)

		err = c.Create(context.TODO(), instance)
		Expect(err).NotTo(HaveOccurred())
		defer c.Delete(context.TODO(), instance)

		createdInstance := &devopsv1alpha1.S2iRun{}
		Eventually(func() error {
			return c.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, createdInstance)
		}, timeout).Should(Succeed())

		instanceUidSlice := strings.Split(string(createdInstance.UID), "-")
		var cmKey = types.NamespacedName{Name: instance.Name + fmt.Sprintf("-%s", instanceUidSlice[len(instanceUidSlice)-1]) + "-configmap", Namespace: "default"}

		//configmap
		cm := &corev1.ConfigMap{}
		Eventually(func() error { return c.Get(context.TODO(), cmKey, cm) }, timeout).
			ShouldNot(Succeed())
	})
})
