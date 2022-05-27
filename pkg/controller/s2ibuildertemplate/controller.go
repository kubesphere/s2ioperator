/*
Copyright 2022 The KubeSphere Authors.

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

package s2ibuildertemplate

import (
	"context"
	"fmt"
	"github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/util/docker"
	"github.com/kubesphere/s2ioperator/pkg/util/reflectutils"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)
import "sigs.k8s.io/controller-runtime/pkg/client"
import "k8s.io/client-go/tools/record"
import "sigs.k8s.io/controller-runtime/pkg/predicate"

// S2iBuilderTemplateReconciler is the reconciler of the S2iBuilderTemplate
type S2iBuilderTemplateReconciler struct {
	client.Client
	recorder record.EventRecorder
}

func (r *S2iBuilderTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	template := &v1alpha1.S2iBuilderTemplate{}
	if err = r.Get(ctx, req.NamespacedName, template); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	if len(template.Spec.ContainerInfo) == 0 {
		r.defaultWarningEvent(template, "spec.containerInfo is missing")
	}

	if template.Spec.DefaultBaseImage == "" {
		r.defaultWarningEvent(template, "spec.defaultBaseImage is missing")
	}
	var builderImages []string
	for _, ImageInfo := range template.Spec.ContainerInfo {
		builderImages = append(builderImages, ImageInfo.BuilderImage)
	}
	if !reflectutils.Contains(template.Spec.DefaultBaseImage, builderImages) {
		r.defaultWarningEvent(template,
			fmt.Sprintf("defaultBaseImage [%s] should in [%v]", template.Spec.DefaultBaseImage, builderImages))
	}

	for i, imageInfo := range template.Spec.ContainerInfo {
		if err := docker.ValidateDockerReference(imageInfo.BuilderImage); err != nil {
			r.defaultWarningEvent(template, fmt.Sprintf("spec.ContainerInfo.[%d]BuilderImage is invalid", i))
		}
	}
	if err := docker.ValidateDockerReference(template.Spec.DefaultBaseImage); err != nil {
		r.defaultWarningEvent(template, "spec.DefaultBaseImage is invalid")
	}
	return
}

func (r *S2iBuilderTemplateReconciler) defaultWarningEvent(template *v1alpha1.S2iBuilderTemplate, message string) {
	r.recorder.Event(template, corev1.EventTypeWarning, "NoRequiredField", message)
}

func (r *S2iBuilderTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	name := "S2iBuilderTemplate"

	r.recorder = mgr.GetEventRecorderFor(name)
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		For(&v1alpha1.S2iBuilderTemplate{}).
		Complete(r)
}
