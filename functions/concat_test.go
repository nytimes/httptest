package functions

import (
	"testing"
)

func TestConcat(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        string
	}{
		{
			map[string]string{},
			[]string{""},
			"",
		},
		{
			map[string]string{},
			[]string{"my string"},
			"my string",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{""},
			"",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header"},
			"123",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header", "my string"},
			"123my string",
		},
		{
			map[string]string{
				"x-previous-header1": "123",
				"x-previous-header2": "456",
			},
			[]string{"x-previous-header1", "my string", "x-previous-header2"},
			"123my string456",
		},
		{
			map[string]string{
				"x-previous-header1": "123",
			},
			[]string{"x-previous-header2", "my string"},
			"x-previous-header2my string",
		},
	}

	for _, tc := range tests {
		actual, _ := Concat(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("Concat(%v, %v): expected %v, actual %v", tc.existingHeaders, tc.args, tc.expected, actual)
		}
	}
}
