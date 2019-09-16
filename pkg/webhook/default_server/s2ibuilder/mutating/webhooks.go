package mutating

import (
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var (
	// Builders contain admission webhook builders
	Builders = map[string]*builder.WebhookBuilder{}
	// HandlerMap contains admission webhook handlers
	HandlerMap = map[string][]admission.Handler{}
)
