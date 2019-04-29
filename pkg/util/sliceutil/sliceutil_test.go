package sliceutil

import (
	"reflect"
	"testing"
)

func TestContainsString(t *testing.T) {
	src := []string{"aa", "bb", "cc"}
	if !ContainsString(src, "bb", nil) {
		t.Errorf("ContainsString didn't find the string as expected")
	}

	modifier := func(s string) string {
		if s == "cc" {
			return "ee"
		}
		return s
	}
	if !ContainsString(src, "ee", modifier) {
		t.Errorf("ContainsString didn't find the string by modifier")
	}
}

func TestRemoveString(t *testing.T) {
	modifier := func(s string) string {
		if s == "ab" {
			return "ee"
		}
		return s
	}
	tests := []struct {
		testName string
		input    []string
		remove   string
		modifier func(s string) string
		want     []string
	}{
		{
			testName: "Nil input slice",
			input:    nil,
			remove:   "",
			modifier: nil,
			want:     nil,
		},
		{
			testName: "Slice doesn't contain the string",
			input:    []string{"a", "ab", "cdef"},
			remove:   "NotPresentInSlice",
			modifier: nil,
			want:     []string{"a", "ab", "cdef"},
		},
		{
			testName: "All strings removed, result is nil",
			input:    []string{"a"},
			remove:   "a",
			modifier: nil,
			want:     nil,
		},
		{
			testName: "No modifier func, one string removed",
			input:    []string{"a", "ab", "cdef"},
			remove:   "ab",
			modifier: nil,
			want:     []string{"a", "cdef"},
		},
		{
			testName: "No modifier func, all(three) strings removed",
			input:    []string{"ab", "a", "ab", "cdef", "ab"},
			remove:   "ab",
			modifier: nil,
			want:     []string{"a", "cdef"},
		},
		{
			testName: "Removed both the string and the modifier func result",
			input:    []string{"a", "cd", "ab", "ee"},
			remove:   "ee",
			modifier: modifier,
			want:     []string{"a", "cd"},
		},
	}
	for _, tt := range tests {
		if got := RemoveString(tt.input, tt.remove, tt.modifier); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%v: RemoveString(%v, %q, %T) = %v WANT %v", tt.testName, tt.input, tt.remove, tt.modifier, got, tt.want)
		}
	}
}
