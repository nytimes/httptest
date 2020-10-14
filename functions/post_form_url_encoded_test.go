package functions

import (
	"errors"
	"net/url"
	"testing"
)

const testURL = "http://httpbin.org/post"

const testResponse1 = `{
"args": {},
"data": "",
"files": {},
"integer": -5,
"decimal": 3.14,
"boolean": true,
"form": {
	"names": [
	  "a",
	  "b"
	]
},
"headers": {
	"Accept": "*/*",
	"Content-Length": "19",
	"Content-Type": "application/x-www-form-urlencoded"
},
"json": null,
"url": "https://httpbin.org/post",
"nested": {
    "object": {
		"data": "123"
	}
  }
}
`

const testResponse2 = `[
"a", "b", "c"
]
`

func TestPostFormURLEncoded(t *testing.T) {
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
		{
			map[string]string{},
			[]string{testURL, "form", "name=abc", "name=xyz"},
			`{"name":["abc","xyz"]}`,
			nil,
		},
		{
			map[string]string{},
			[]string{testURL, "form.name", "name=abc", "name=xyz"},
			`["abc","xyz"]`,
			nil,
		},
		{
			map[string]string{},
			[]string{testURL, "form.name.1", "name=abc", "name=xyz"},
			`xyz`,
			nil,
		},
	}

	for _, tc := range tests {
		actual, err := PostFormURLEncoded(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("PostFormURLEncoded(%v, %v): expected %v, actual %v", tc.existingHeaders, tc.args, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("PostFormURLEncoded(%v, %v): expected %v, got: %v", tc.existingHeaders, tc.args, tc.err, err)
		}
	}
}

func TestArgsToRequestBody(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        url.Values
	}{
		{
			map[string]string{},
			[]string{"foo=bar", "abc=123", "xyz=456"},
			map[string][]string{
				"foo": []string{"bar"},
				"abc": []string{"123"},
				"xyz": []string{"456"},
			},
		},
		{
			map[string]string{},
			[]string{"foo=bar"},
			map[string][]string{
				"foo": []string{"bar"},
			},
		},
		{
			map[string]string{
				"x-previous-header1": "foo=bar",
				"x-previous-header2": "abc=123",
			},
			[]string{"foo=bar", "x-previous-header2", "xyz=456"},
			map[string][]string{
				"foo": []string{"bar"},
				"abc": []string{"123"},
				"xyz": []string{"456"},
			},
		},
	}

	for _, tc := range tests {
		actual := argsToRequestBody(tc.existingHeaders, tc.args)
		if actual.Encode() != tc.expected.Encode() {
			t.Errorf("argsToRequestBody(%v, %v): expected %v, actual %v", tc.existingHeaders, tc.args, tc.expected, actual)
		}
	}
}

func TestRetrieveElementTree(t *testing.T) {
	var tests = []struct {
		json     string
		element  string
		expected string
		err      error
	}{
		{
			testResponse1,
			"url",
			"https://httpbin.org/post",
			nil,
		},
		{
			testResponse1,
			"integer",
			"-5",
			nil,
		},
		{
			testResponse1,
			"decimal",
			"3.14",
			nil,
		},
		{
			testResponse1,
			"boolean",
			"true",
			nil,
		},
		{
			testResponse1,
			"json",
			"null",
			nil,
		},
		{
			testResponse1,
			"data",
			"",
			nil,
		},
		{
			testResponse1,
			"files",
			"{}",
			nil,
		},
		{
			testResponse1,
			"headers",
			`{"Accept":"*/*","Content-Length":"19","Content-Type":"application/x-www-form-urlencoded"}`,
			nil,
		},
		{
			testResponse1,
			"headers.Content-Length",
			"19",
			nil,
		},
		{
			testResponse1,
			"form.names",
			`["a","b"]`,
			nil,
		},
		{
			testResponse1,
			"form.names.1",
			"b",
			nil,
		},
		{
			testResponse1,
			"nested.object.data",
			"123",
			nil,
		},
		{
			testResponse2,
			"1",
			"b",
			nil,
		},
		{
			testResponse1,
			"headers.User-Agent",
			"",
			errors.New("error: element User-Agent not found in response body"),
		},
		{
			testResponse1,
			"form.names.2",
			"",
			errors.New("JSON array index out of bounds: 2"),
		},
		{
			testResponse1,
			"form.names.foo",
			"",
			errors.New("invalid index for JSON array: foo"),
		},
		{
			testResponse1,
			"nested.object.bogus",
			"",
			errors.New("error: element bogus not found in response body"),
		},
		{
			testResponse2,
			"3",
			"",
			errors.New("JSON array index out of bounds: 3"),
		},
		{
			testResponse1,
			"form.names.-1",
			"",
			errors.New("JSON array index out of bounds: -1"),
		},
	}

	for _, tc := range tests {
		actual, err := retrieveElementTree([]byte(tc.json), tc.element)
		if string(actual) != tc.expected {
			t.Errorf("retrieveElementTree(%v, %v): expected %v, actual %v", tc.json, tc.element, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("retrieveElementTree(%v, %v): expected %v, got: %v", tc.json, tc.element, tc.err, err)
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
		actual, err := jsonToBytes(tc.rawJSON)
		if string(actual) != tc.expected {
			t.Errorf("jsonToBytes(%v): expected %v, actual %v", tc.rawJSON, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("jsonToBytes(%v): expected %v, got: %v", tc.rawJSON, tc.err, err)
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
