package reflectutils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestContains(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Contains Suite")
}

var _ = Describe("Testing reflectutils Contains fuc string slice", func() {
	It("should run without err", func() {

		inSliceValues := []string{"a", "b", "c", "d"}

		outOfSliceValues := []string{"cc", "dd", "dfsfds;;", "//.,,"}

		stringSlices := [][]string{
			{"a", "b"},
			{"b", "c"},
			{"aaaaaddd", "cde", "c"},
			{"!!!", "dsfdsfs", "d"},
		}
		for i, test := range stringSlices {
			Expect(Contains(inSliceValues[i], test), BeTrue())
		}

		for i, test := range stringSlices {
			Expect(Contains(outOfSliceValues[i], test)).To(BeFalse())
		}
	})
})
