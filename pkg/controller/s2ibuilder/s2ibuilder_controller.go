/*
Copyright 2019 The Kubesphere Authors.

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

package s2ibuilder

import (
	"context"
	"reflect"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("s2ibuilder-controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new S2iBuilder Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileS2iBuilder{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("s2ibuilder-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to S2iBuilder
	err = c.Watch(&source.Kind{Type: &devopsv1alpha1.S2iBuilder{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//watch s2irun
	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			run := a.Object.(*devopsv1alpha1.S2iRun)
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      run.Spec.BuilderName,
					Namespace: a.Meta.GetNamespace(),
				}},
			}
		})

	// 'UpdateFunc' and 'CreateFunc' used to judge if a event about the object is
	// what we want. If that is true, the event will be processed by the reconciler.
	p := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			run := e.ObjectOld.(*devopsv1alpha1.S2iRun)
			if run.Spec.BuilderName == "" {
				return false
			}
			return e.ObjectOld != e.ObjectNew
		},
		CreateFunc: func(e event.CreateEvent) bool {
			run := e.Object.(*devopsv1alpha1.S2iRun)
			if run.Spec.BuilderName == "" {
				return false
			}
			return true
		},
	}
	err = c.Watch(&source.Kind{Type: &devopsv1alpha1.S2iRun{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: mapFn,
	}, p)

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileS2iBuilder{}

// ReconcileS2iBuilder reconciles a S2iBuilder object
type ReconcileS2iBuilder struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a S2iBuilder object and makes changes based on the state read
// and what is in the S2iBuilder.Spec
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=s2ibuilders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=s2ibuilders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
func (r *ReconcileS2iBuilder) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the S2iBuilder instance
	log.Info("Reconciler of s2ibuilder called", "NamespaceName", request.NamespacedName)
	instance := &devopsv1alpha1.S2iBuilder{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	origin := instance.DeepCopy()
	runList := new(devopsv1alpha1.S2iRunList)
	err = r.Client.List(context.TODO(), client.InNamespace(instance.Namespace), runList)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	instance.Status.RunCount = 0
	last := new(metav1.Time)
	for _, item := range runList.Items {
		if item.Spec.BuilderName == instance.Name {
			instance.Status.RunCount++
			if item.Status.StartTime != nil && item.Status.StartTime.After(last.Time) {
				*last = *(item.Status.StartTime)
				instance.Status.LastRunState = item.Status.RunState
				if instance.Status.LastRunName == nil {
					instance.Status.LastRunName = new(string) //should use defaulting instead of creating here
				}
				*(instance.Status.LastRunName) = item.Name
				if instance.Status.LastRunStartTime == nil {
					instance.Status.LastRunStartTime = new(metav1.Time)
				}
				*(instance.Status.LastRunStartTime) = *item.Status.StartTime
			}
		}
	}
	if instance.Status.RunCount == 0 {
		instance.Status.LastRunName = nil
		instance.Status.LastRunState = ""
		instance.Status.LastRunStartTime = nil
	}
	if !reflect.DeepEqual(instance.Status, origin.Status) {
		if err := r.Status().Update(context.Background(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
