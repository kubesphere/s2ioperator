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

package s2irun

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	loghandler "github.com/kubesphere/s2ioperator/pkg/handler/log"
	"github.com/kubesphere/s2ioperator/pkg/util/reflectutils"
	"k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("s2irun-controller")

const (
	S2iRunBuilderLabel = "labels.devops.kubesphere.io/builder-name"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new S2iRun Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileS2iRun{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("s2irun-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to S2iRun
	err = c.Watch(&source.Kind{Type: &devopsv1alpha1.S2iRun{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devopsv1alpha1.S2iRun{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &devopsv1alpha1.S2iRun{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileS2iRun{}

// ReconcileS2iRun reconciles a S2iRun object
type ReconcileS2iRun struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a S2iRun object and makes changes based on the state read
// and what is in the S2iRun.Spec

// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=s2iruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=s2iruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=devops.kubesphere.io,resources=s2ibuildertemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=extensions,resources=deployments,verbs=get;list;watch;create;update;patch
func (r *ReconcileS2iRun) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the S2iRun instance
	log.Info("Reconciler of s2irun called", "Name", request.Name)
	instance := &devopsv1alpha1.S2iRun{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if k8serror.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	origin := instance.DeepCopy()
	//configmap setup
	builder := &devopsv1alpha1.S2iBuilder{}
	if err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.BuilderName, Namespace: instance.Namespace}, builder); err != nil {
		if k8serror.IsNotFound(err) {
			log.Info("Waiting for creating s2ibuilder", "Name", instance.Spec.BuilderName)
			return reconcile.Result{RequeueAfter: time.Second * 15}, nil
		}
		return reconcile.Result{}, err
	}
	if instance.Labels == nil {
		instance.Labels = make(map[string]string)
	}
	if v, ok := instance.Labels[S2iRunBuilderLabel]; !ok || v != builder.Name {
		instance.Labels[S2iRunBuilderLabel] = builder.Name
		err = r.Update(context.TODO(), instance)
		if err != nil {
			log.Error(nil, "Failed to add labels to s2irun")
			return reconcile.Result{}, err
		}
	}
	configmap, err := r.NewConfigMap(instance, *builder.Spec.Config, builder.Spec.FromTemplate)
	if err != nil {
		log.Error(err, "Failed to initialize a configmap")
		return reconcile.Result{}, err
	}
	foundcm := &corev1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: configmap.Name, Namespace: configmap.Namespace}, foundcm)
	if err != nil && k8serror.IsNotFound(err) {
		log.Info("Creating ConfigMap", "Namespace", configmap.Namespace, "name", configmap.Name)
		if err := controllerutil.SetControllerReference(instance, configmap, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.Create(context.TODO(), configmap)
		if err != nil {
			log.Error(err, "Create configmap failed", "Namespace", configmap.Namespace, "name", configmap.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		if !reflect.DeepEqual(configmap.Data, foundcm.Data) {
			foundcm.Data = configmap.Data
			log.Info("Updating job config", "Namespace", foundcm.Namespace, "Name", foundcm.Name)
			err = r.Update(context.TODO(), foundcm)
			if err != nil {
				log.Error(err, "Failed to updating job config", "Namespace", foundcm.Namespace, "Name", foundcm.Name)
				return reconcile.Result{}, err
			}
		}
	}
	//job set up
	job, err := r.GenerateNewJob(instance)
	if err != nil {
		log.Error(err, "Failed to initialize a job")
		return reconcile.Result{}, err
	}
	found := &batchv1.Job{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && k8serror.IsNotFound(err) {
		log.Info("Creating Job", "Namespace", job.Namespace, "Name", job.Name)
		if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.Create(context.TODO(), job)
		if err != nil {
			log.Error(err, "Failed to create Job", "Namespace", job.Namespace, "Name", job.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		instance.Status.KubernetesJobName = found.Name
		instance.Status.StartTime = found.Status.StartTime
		if found.Status.Active == 1 {
			log.Info("Job is running", "start time", found.Status.StartTime)
			instance.Status.RunState = devopsv1alpha1.Running
		} else if found.Status.Failed == 1 {
			log.Info("Job failed")
			instance.Status.RunState = devopsv1alpha1.Failed
			instance.Status.CompletionTime = found.Status.CompletionTime
		} else if found.Status.Succeeded == 1 {
			log.Info("Job completed", "time", found.Status.CompletionTime)
			instance.Status.RunState = devopsv1alpha1.Successful
			instance.Status.CompletionTime = found.Status.CompletionTime
			logURL, err := r.GetLogURL(found)
			if err != nil {
				return reconcile.Result{}, err
			}
			instance.Status.LogURL = logURL

		} else {
			instance.Status.RunState = devopsv1alpha1.Unknown
		}
	}
	if !reflect.DeepEqual(instance.Status, origin.Status) {
		err = r.Status().Update(context.Background(), instance)
		if err != nil {
			log.Error(nil, "Failed to update s2irun status", "Namespace", instance.Namespace, "Name", instance.Name)
			return reconcile.Result{}, err
		}
	} else if instance.Status.RunState == devopsv1alpha1.Successful {
		err := r.ScaleWorkLoads(instance, builder)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileS2iRun) GetLogURL(job *batchv1.Job) (string, error) {
	pods := &corev1.PodList{}
	listOption := &client.ListOptions{}
	listOption.SetLabelSelector("job-name=" + job.Name)
	listOption.InNamespace(job.Namespace)
	err := r.List(context.TODO(), listOption, pods)
	if err != nil {
		log.Error(nil, "Error in get pod of job")
		return "", nil
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("cannot find any pod of the job %s", job.Name)
	}
	return loghandler.GetKubesphereLogger().GetURLOfPodLog(pods.Items[0].Namespace, pods.Items[0].Name)
}

func GetNewImangeName(instance *devopsv1alpha1.S2iRun, config devopsv1alpha1.S2iConfig) string {
	if instance.Spec.NewTag != "" {
		return config.ImageName + ":" + instance.Spec.NewTag
	} else {
		return config.ImageName + ":" + config.Tag
	}
}

// ScaleWorkLoads will auto scale workloads define in s2ibuilder's annotations
func (r *ReconcileS2iRun) ScaleWorkLoads(instance *devopsv1alpha1.S2iRun, builder *devopsv1alpha1.S2iBuilder) error {
	if annotation, ok := builder.Annotations[devopsv1alpha1.AutoScaleAnnotations]; ok {
		log.Info("Start AutoScale Workloads")
		origin := instance.DeepCopy()
		s2iAutoScale := make([]devopsv1alpha1.S2iAutoScale, 0)
		completedScaleWorkloads := make([]devopsv1alpha1.S2iAutoScale, 0)
		if err := json.Unmarshal([]byte(annotation), &s2iAutoScale); err != nil {
			return err
		}
		errs := make([]error, 0)
		if completedScaleAnnotations, ok := instance.Annotations[devopsv1alpha1.S2irCompletedScaleAnnotations]; ok {
			if err := json.Unmarshal([]byte(completedScaleAnnotations), &completedScaleWorkloads); err != nil {
				return err
			}
		}
		for _, scale := range s2iAutoScale {
			hasScaled := false
			for _, completedScale := range completedScaleWorkloads {
				if reflect.DeepEqual(scale, completedScale) {
					hasScaled = true
					break
				}
			}
			if hasScaled {
				continue
			}
			switch scale.Kind {

			case devopsv1alpha1.KindDeployment:
				deploy := &v1.Deployment{}
				err := r.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: scale.Name}, deploy)
				if err != nil && k8serror.IsNotFound(err) {
					errs = append(errs, err)
					continue
				} else if err != nil {
					return err
				}

				log.Info("Autoscale Deployment", "ns", instance.Namespace, "statefulSet", deploy.Name)
				newImageName := GetNewImangeName(instance, *builder.Spec.Config)
				// if only one container update containers image config
				if len(deploy.Spec.Template.Spec.Containers) == 1 {
					if deploy.Spec.Template.Spec.Containers[0].Image == newImageName {
						deploy.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullAlways
					} else {
						deploy.Spec.Template.Spec.Containers[0].Image = newImageName
					}
				} else {
					for _, container := range deploy.Spec.Template.Spec.Containers {
						if reflectutils.Contains(container.Name, scale.Containers) {
							if container.Image == newImageName {
								container.ImagePullPolicy = corev1.PullAlways
							} else {
								container.Image = newImageName
							}
						}
					}
				}
				if deploy.Spec.Template.Labels == nil {
					deploy.Spec.Template.Labels = make(map[string]string)
				}

				deploy.Spec.Template.Labels[devopsv1alpha1.WorkloadLatestS2iRunTemplateLabel] = instance.Name

				log.Info("Update deployment", "ns", deploy.Namespace, "statefulSet", deploy.Name)
				err = r.Update(context.TODO(), deploy)
				if err != nil && k8serror.IsNotFound(err) {
					errs = append(errs, err)
					continue
				} else if err != nil {
					return err
				}
				completedScaleWorkloads = append(completedScaleWorkloads, scale)

			case devopsv1alpha1.KindStatefulSet:
				statefulSet := &v1.StatefulSet{}
				err := r.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: scale.Name}, statefulSet)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				log.Info("Autoscale StatefulSet", "ns", instance.Namespace, "statefulSet", statefulSet.Name)
				newImageName := GetNewImangeName(instance, *builder.Spec.Config)
				if len(statefulSet.Spec.Template.Spec.Containers) == 1 {
					if statefulSet.Spec.Template.Spec.Containers[0].Image == newImageName {
						statefulSet.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullAlways
					} else {
						statefulSet.Spec.Template.Spec.Containers[0].Image = newImageName
					}
				} else {
					for _, container := range statefulSet.Spec.Template.Spec.Containers {
						if reflectutils.Contains(container.Name, scale.Containers) {
							if container.Image == newImageName {
								container.ImagePullPolicy = corev1.PullAlways
							} else {
								container.Image = newImageName
							}
						}
					}
				}
				if statefulSet.Spec.Template.Labels == nil {
					statefulSet.Spec.Template.Labels = make(map[string]string)
				}

				statefulSet.Spec.Template.Labels[devopsv1alpha1.WorkloadLatestS2iRunTemplateLabel] = instance.Name

				log.Info("Update statefulSet", "ns", statefulSet.Namespace, "statefulSet", statefulSet.Name)
				err = r.Update(context.TODO(), statefulSet)
				if err != nil && k8serror.IsNotFound(err) {
					errs = append(errs, err)
					continue
				} else if err != nil {
					return err
				}
				completedScaleWorkloads = append(completedScaleWorkloads, scale)
			default:
				errs = append(errs, fmt.Errorf("unsupport workload Kind [%s], name [%s]", scale.Kind, scale.Name))
			}
		}
		if completedScaleAnnotation, err := json.Marshal(completedScaleWorkloads); err != nil {
			return err
		} else {
			instance.Annotations[devopsv1alpha1.S2irCompletedScaleAnnotations] = string(completedScaleAnnotation)
		}
		if !reflect.DeepEqual(origin, instance) {
			if err := r.Update(context.TODO(), instance); err != nil {
				return err
			}
		}
		if len(errs) != 0 {
			return errors.NewAggregate(errs)
		}
	}
	return nil
}
