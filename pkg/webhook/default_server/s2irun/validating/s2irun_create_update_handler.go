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

package validating

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/errors"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	types2 "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-s2irun"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &S2iRunCreateUpdateHandler{})
}

// S2iRunCreateUpdateHandler handles S2iBuilder
type S2iRunCreateUpdateHandler struct {
	Client client.Client
	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *S2iRunCreateUpdateHandler) validatingS2iRunFn(ctx context.Context, obj *devopsv1alpha1.S2iRun) (bool, string, error) {
	origin := &devopsv1alpha1.S2iRun{}

	builder := &devopsv1alpha1.S2iBuilder{}

	err := h.Client.Get(context.TODO(), types2.NamespacedName{Namespace: obj.Namespace, Name: obj.Spec.BuilderName}, builder)
	if err != nil && !k8serror.IsNotFound(err) {
		return false, "validate failed", errors.NewFieldInvalidValueWithReason("no", "could not call k8s api")
	}
	if !k8serror.IsNotFound(err) {
		if obj.Spec.NewSourceURL != "" && !builder.Spec.Config.IsBinaryURL {
			return false, "validate failed", errors.NewFieldInvalidValueWithReason("newSourceURL", "only b2i could set newSourceURL")
		}
	}

	err = h.Client.Get(context.TODO(), types2.NamespacedName{Namespace: obj.Namespace, Name: obj.Name}, origin)
	if !k8serror.IsNotFound(err) && origin.Status.RunState != "" && !reflect.DeepEqual(origin.Spec, obj.Spec) {
		return false, "validate failed", errors.NewFieldInvalidValueWithReason("spec", "should not change s2i run spec when job started")
	}

	if obj.Spec.NewTag != "" {
		validateImageName := fmt.Sprintf("validate:%s", obj.Spec.NewTag)
		if err := validateDockerReference(validateImageName); err != nil {
			return false, "", err
		}
	}

	return true, "", nil
}

var _ admission.Handler = &S2iRunCreateUpdateHandler{}

// Handle handles admission requests.
func (h *S2iRunCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &devopsv1alpha1.S2iRun{}
	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatingS2iRunFn(ctx, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Decoder = &S2iRunCreateUpdateHandler{}

// InjectClient injects the client into the S2iBuilderCreateUpdateHandler
func (h *S2iRunCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

// InjectDecoder injects the decoder into the S2iRunCreateUpdateHandler
func (h *S2iRunCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
