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

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-s2ibuilder"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &S2iBuilderCreateUpdateHandler{})
}

// S2iBuilderCreateUpdateHandler handles S2iBuilder
type S2iBuilderCreateUpdateHandler struct {
	Client client.Client
	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *S2iBuilderCreateUpdateHandler) validatingS2iBuilderFn(ctx context.Context, obj *devopsv1alpha1.S2iBuilder) (bool, string, error) {
	// TODO(user): implement your admission logic
	fromTemplate := false
	if obj.Spec.FromTemplate != nil {
		t := &devopsv1alpha1.S2iBuilderTemplate{}
		err := h.Client.Get(context.TODO(), k8stypes.NamespacedName{Name: obj.Spec.FromTemplate.Name}, t)
		if err != nil {
			if k8serror.IsNotFound(err) {
				return false, "validate failed", fmt.Errorf("Template not found, pls check the template name  [%s] or create a template", obj.Spec.FromTemplate.Name)
			}
			return false, "Unhandle error", err
		}
		fromTemplate = true
	}

	if errs := ValidateConfig(obj.Spec.Config, fromTemplate); len(errs) == 0 {
		return true, "allowed to be admitted", nil
	} else {
		return false, "validate failed", errors.NewAggregate(errs)
	}
}

var _ admission.Handler = &S2iBuilderCreateUpdateHandler{}

// Handle handles admission requests.
func (h *S2iBuilderCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &devopsv1alpha1.S2iBuilder{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatingS2iBuilderFn(ctx, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Client = &S2iBuilderCreateUpdateHandler{}

// InjectClient injects the client into the S2iBuilderCreateUpdateHandler
func (h *S2iBuilderCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

var _ inject.Decoder = &S2iBuilderCreateUpdateHandler{}

// InjectDecoder injects the decoder into the S2iBuilderCreateUpdateHandler
func (h *S2iBuilderCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
