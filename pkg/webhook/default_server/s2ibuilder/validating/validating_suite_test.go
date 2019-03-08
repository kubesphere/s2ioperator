package validating_test

import (
	"strings"
	"testing"

	api "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/webhook/default_server/s2ibuilder/validating"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidating(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validating Suite")
}

var _ = Describe("Valiating test", func() {
	var (
		User   []api.Parameter
		Define []api.Parameter
	)
	BeforeEach(func() {
		User = []api.Parameter{
			api.Parameter{
				Key:   "Key1",
				Value: "Value1",
			}, api.Parameter{
				Key:   "Key2",
				Value: "Value2",
			}, api.Parameter{
				Key:   "Key3",
				Value: "Value3",
			}}

		Define = []api.Parameter{api.Parameter{
			Key:      "Key1",
			Required: true,
		},
			api.Parameter{
				Key:      "Key2",
				Required: true,
			},
			api.Parameter{
				Key: "Key3",
			}}
	})
	It("Should pass when all field set", func() {
		errors := validating.ValidateParameter(User, Define)
		Expect(errors).To(HaveLen(0))
	})
	It("Should not pass when missing Key1", func() {
		User[0].Value = ""
		errors := validating.ValidateParameter(User, Define)
		Expect(errors).To(HaveLen(1))
		Expect(strings.Contains(errors[0].Error(), User[0].Key)).To(BeTrue())
	})
	It("Should Pass when key3 is ignored", func() {
		User = User[:2]
		errors := validating.ValidateParameter(User, Define)
		Expect(errors).To(HaveLen(0))
	})
})
