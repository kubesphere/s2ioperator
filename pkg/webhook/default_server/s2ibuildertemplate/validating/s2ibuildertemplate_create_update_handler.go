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
	"github.com/kubesphere/s2ioperator/pkg/errors"
	"github.com/kubesphere/s2ioperator/pkg/util/reflectutils"
	"net/http"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-s2ibuildertemplate"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &S2iBuilderTemplateCreateUpdateHandler{})
}

// S2iBuilderTemplateCreateUpdateHandler handles S2iBuilder
type S2iBuilderTemplateCreateUpdateHandler struct {
	Client client.Client
	// Decoder decodes objects
	Decoder types.Decoder
}

func (h *S2iBuilderTemplateCreateUpdateHandler) validatingS2iBuilderTemplateFn(ctx context.Context, obj *devopsv1alpha1.S2iBuilderTemplate) (bool, string, error) {

	if len(obj.Spec.BaseImages) == 0 {
		return false, "validate failed", errors.NewFieldRequired("baseImages")
	}

	if obj.Spec.DefaultBaseImage == "" {
		return false, "validate failed", errors.NewFieldRequired("defaultBaseImage")
	}

	if !reflectutils.Contains(obj.Spec.DefaultBaseImage, obj.Spec.BaseImages) {
		return false, "validate failed", errors.NewFieldInvalidValueWithReason("defaultBaseImage",
			fmt.Sprintf("defaultBaseImage [%s] should in [%v]", obj.Spec.DefaultBaseImage, obj.Spec.BaseImages))
	}

	for _, baseImage := range obj.Spec.BaseImages {
		if err := validateDockerReference(baseImage); err != nil {
			return false, "validate failed", errors.NewFieldInvalidValueWithReason("builderImage", err.Error())
		}
	}
	if err := validateDockerReference(obj.Spec.DefaultBaseImage); err != nil {
		return false, "validate failed", errors.NewFieldInvalidValueWithReason("defaultBaseImage", err.Error())
	}
	return true, "", nil
}

var _ admission.Handler = &S2iBuilderTemplateCreateUpdateHandler{}

// Handle handles admission requests.
func (h *S2iBuilderTemplateCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &devopsv1alpha1.S2iBuilderTemplate{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	allowed, reason, err := h.validatingS2iBuilderTemplateFn(ctx, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

var _ inject.Decoder = &S2iBuilderTemplateCreateUpdateHandler{}

// InjectDecoder injects the decoder into the S2iBuilderTemplateCreateUpdateHandler
func (h *S2iBuilderTemplateCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
