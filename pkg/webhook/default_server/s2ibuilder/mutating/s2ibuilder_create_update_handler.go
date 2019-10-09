package mutating

import (
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

const (
	DefaultRevisionId = "master"
	DefaultTag        = "latest"
)

func init() {
	webhookName := "mutating-create-update-s2ibuilder"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &S2iBuilderCreateHandler{})
}

// S2iBuilderCreateHandler handles S2iBuilder
type S2iBuilderCreateHandler struct {
	Decoder types.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &S2iBuilderCreateHandler{}

// S2iBuilderCreateHandler adds an default status info to S2iBuilder
func (h *S2iBuilderCreateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	s2ibuilder := &devopsv1alpha1.S2iBuilder{}
	err := h.Decoder.Decode(req, s2ibuilder)

	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	s2ib := s2ibuilder.DeepCopy()

	err = h.mutatingS2iBuilderFn(ctx, s2ib)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(s2ibuilder, s2ib)
}

func (h *S2iBuilderCreateHandler) mutatingS2iBuilderFn(ctx context.Context, obj *devopsv1alpha1.S2iBuilder) error {

	if obj.Spec.Config.RevisionId == "" {
		obj.Spec.Config.RevisionId = DefaultRevisionId
	}

	if obj.Spec.Config.Tag == "" {
		obj.Spec.Config.Tag = DefaultTag
	}

	return nil
}

var _ inject.Decoder = &S2iBuilderCreateHandler{}

// InjectDecoder injects the decoder.
func (h *S2iBuilderCreateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
