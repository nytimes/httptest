package functions

import (
	"errors"
	"testing"
)

const testURL = "http://httpbin.org/post"

func TestPostData(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        string
		err             error
	}{
		{
			map[string]string{},
			[]string{},
			"",
			errURLMissing,
		},
		{
			map[string]string{},
			[]string{"", "foo", "bar"},
			"",
			errURLEmpty,
		},
		{
			map[string]string{},
			[]string{testURL, "url"},
			testURL,
			nil,
		},
		{
			map[string]string{},
			[]string{testURL, "foo"},
			"",
			errors.New("error: element foo not found in response body"),
		},
		{
			map[string]string{},
			[]string{testURL, "form", "data1=123", "data2=456", "data3=789"},
			`{"data1":"123","data2":"456","data3":"789"}`,
			nil,
		},
		{
			map[string]string{
				"x-previous-header1": "foo=bar",
				"x-previous-header2": "abc=xyz",
			},
			[]string{testURL, "form", "data1=123", "x-previous-header2", "data2=456"},
			`{"abc":"xyz","data1":"123","data2":"456"}`, // this API will return the keys in alphabetical order
			nil,
		},
	}

	for _, tc := range tests {
		actual, err := PostData(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("PostData(%v, %v): expected %v, actual %v", tc.existingHeaders, tc.args, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("PostData(%v, %v): expected %v, got: %v", tc.existingHeaders, tc.args, tc.err, err)
		}
	}
}

func TestArgsToURLResponseElement(t *testing.T) {
	var tests = []struct {
		args                    []string
		expectedURL             string
		expectedResponseElement string
		err                     error
	}{
		{
			[]string{"foo", "bar"},
			"foo",
			"bar",
			nil,
		},
		{
			[]string{"foo", "bar", "abc", "xyz"},
			"foo",
			"bar",
			nil,
		},
		{
			[]string{"", "bar", "abc", "xyz"},
			"",
			"",
			errURLEmpty,
		},
		{
			[]string{"foo", "", "abc", "xyz"},
			"foo",
			"",
			nil,
		},
	}

	for _, tc := range tests {
		actualURL, actualResponseElement, err := argsToURLResponseElement(tc.args)
		if actualURL != tc.expectedURL {
			t.Errorf("argsToURLResponseElement(%v): expected URL %v, actual URL %v", tc.args, tc.expectedURL, actualURL)
		}
		if actualResponseElement != tc.expectedResponseElement {
			t.Errorf("argsToURLResponseElement(%v): expected response element %v, actual response element %v", tc.args, tc.expectedResponseElement, actualResponseElement)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("argsToURLResponseElement(%v): expected %v, got: %v", tc.args, tc.err, err)
		}
	}
}

func TestArgsToRequestBody(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        string
	}{
		{
			map[string]string{},
			[]string{"foo=bar", "abc=123", "xyz=456"},
			"foo=bar&abc=123&xyz=456",
		},
		{
			map[string]string{},
			[]string{"foo=bar"},
			"foo=bar",
		},
		{
			map[string]string{
				"x-previous-header1": "foo=bar",
				"x-previous-header2": "abc=123",
			},
			[]string{"foo=bar", "x-previous-header2", "xyz=456"},
			"foo=bar&abc=123&xyz=456",
		},
	}

	for _, tc := range tests {
		actual := argsToRequestBody(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("argsToRequestBody(%v, %v): expected %v, actual %v", tc.existingHeaders, tc.args, tc.expected, actual)
		}
	}
}

func TestIsJSONObject(t *testing.T) {
	object := `{"foo":"bar","abc":"123"}`
	notObject1 := `foo`
	notObject2 := `[{"foo":"bar"},{"abc":"123"}]`

	var tests = []struct {
		str      []byte
		expected bool
	}{
		{
			[]byte(object),
			true,
		},
		{
			[]byte(notObject1),
			false,
		},
		{
			[]byte(notObject2),
			false,
		},
	}

	for _, tc := range tests {
		actual := isJSONObject([]byte(tc.str))
		if actual != tc.expected {
			t.Errorf("isJSONObject(%v): expected %v, actual %v", tc.str, tc.expected, actual)
		}
	}
}

func TestIsJSONArray(t *testing.T) {
	array := `[{"foo":"bar"},{"abc":"123"}]`
	notArray1 := `3.14`
	notArray2 := `{"foo":"bar","abc":"123"}`

	var tests = []struct {
		str      []byte
		expected bool
	}{
		{
			[]byte(array),
			true,
		},
		{
			[]byte(notArray1),
			false,
		},
		{
			[]byte(notArray2),
			false,
		},
	}

	for _, tc := range tests {
		actual := isJSONArray([]byte(tc.str))
		if actual != tc.expected {
			t.Errorf("isJSONArray(%v): expected %v, actual %v", tc.str, tc.expected, actual)
		}
	}
}

func TestJsonToString(t *testing.T) {
	object := `{"foo":"bar","abc":"123"}`
	array := `[{"foo":"bar"},{"abc":"123"}]`
	stringPrimitive := `"foo"`
	numberPrimitive := `3.14`
	booleanPrimitive := `true`

	var tests = []struct {
		rawJSON  []byte
		expected string
		err      error
	}{
		{
			[]byte(object),
			object,
			nil,
		},
		{
			[]byte(array),
			array,
			nil,
		},
		{
			[]byte(stringPrimitive),
			stringPrimitive[1 : len(stringPrimitive)-1],
			nil,
		},
		{
			[]byte(numberPrimitive),
			numberPrimitive,
			nil,
		},
		{
			[]byte(booleanPrimitive),
			booleanPrimitive,
			nil,
		},
	}

	for _, tc := range tests {
		actual, err := jsonToString(tc.rawJSON)
		if actual != tc.expected {
			t.Errorf("jsonToString(%v): expected %v, actual %v", tc.rawJSON, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("jsonToString(%v): expected %v, got: %v", tc.rawJSON, tc.err, err)
		}
	}
}
