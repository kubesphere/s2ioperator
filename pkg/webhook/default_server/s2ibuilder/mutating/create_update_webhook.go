package mutating

import (
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

var logger = log.Log.WithName("s2ibuilder-mutate")

func init() {
	builderName := "mutating-create-update-s2ibuilder"
	Builders[builderName] = builder.
		NewWebhookBuilder().
		Name(builderName+".kubesphere.io").
		Path("/"+builderName).
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		FailurePolicy(admissionregistrationv1beta1.Fail).
		ForType(&devopsv1alpha1.S2iBuilder{})
}
