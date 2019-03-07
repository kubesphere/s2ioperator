package reflectutils

import "testing"

func TestStringSliceContains(t *testing.T) {

	inSliceValues := []string{"a", "b", "c", "d"}

	outOfSliceValues := []string{"cc", "dd", "dfsfds;;", "//.,,"}

	stringSlices := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"aaaaaddd", "cde", "c"},
		{"!!!", "dsfdsfs", "d"},
	}
	for i, test := range stringSlices {
		if !Contains(inSliceValues[i], test) {
			t.Fatalf("%s should in %+v", inSliceValues[i], test)
		}
	}

	for i, test := range stringSlices {
		if Contains(outOfSliceValues[i], test) {
			t.Fatalf("%s should out of %+v", inSliceValues[i], test)
		}
	}

}
