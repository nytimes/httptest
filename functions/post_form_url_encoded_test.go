package functions

import (
	"errors"
	"net/url"
	"testing"
)

const testURL = "http://httpbin.org/post"

const testComplexResponse = `{
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

const testJSONArray = `[
"a", "b", "c"
]
`

const testBadResponse = `["a", "b", "c"`

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
			nil,
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

func TestRetrieveElement(t *testing.T) {
	var tests = []struct {
		json     string
		element  string
		expected string
		err      error
	}{
		{
			testComplexResponse,
			"url",
			"https://httpbin.org/post",
			nil,
		},
		{
			testComplexResponse,
			"integer",
			"-5",
			nil,
		},
		{
			testComplexResponse,
			"decimal",
			"3.14",
			nil,
		},
		{
			testComplexResponse,
			"boolean",
			"true",
			nil,
		},
		{
			testComplexResponse,
			"json",
			"",
			nil,
		},
		{
			testComplexResponse,
			"data",
			"",
			nil,
		},
		{
			testComplexResponse,
			"files",
			"{}",
			nil,
		},
		{
			testComplexResponse,
			"headers",
			`{"Accept":"*/*","Content-Length":"19","Content-Type":"application/x-www-form-urlencoded"}`,
			nil,
		},
		{
			testComplexResponse,
			"headers.Content-Length",
			"19",
			nil,
		},
		{
			testComplexResponse,
			"form.names",
			`["a","b"]`,
			nil,
		},
		{
			testComplexResponse,
			"form.names.1",
			"b",
			nil,
		},
		{
			testComplexResponse,
			"nested.object.data",
			"123",
			nil,
		},
		{
			testJSONArray,
			"1",
			"b",
			nil,
		},
		// non-existent values
		{
			testComplexResponse,
			"headers.User-Agent",
			"",
			nil,
		},
		{
			testComplexResponse,
			"form.names.2",
			"",
			nil,
		},
		{
			testComplexResponse,
			"form.names.foo",
			"",
			nil,
		},
		{
			testComplexResponse,
			"nested.object.bogus",
			"",
			nil,
		},
		{
			testJSONArray,
			"3",
			"",
			nil,
		},
		// negative index
		{
			testComplexResponse,
			"form.names.-1",
			"",
			nil,
		},
		{
			testBadResponse,
			"",
			"",
			errors.New("invalid JSON"),
		},
		{
			testJSONArray,
			"",
			`["a","b","c"]`,
			nil,
		},
	}

	for _, tc := range tests {
		actual, err := retrieveElement([]byte(tc.json), tc.element)
		if string(actual) != tc.expected {
			t.Errorf("retrieveElement(%v, %v): expected %v, actual %v", tc.json, tc.element, tc.expected, actual)
		}
		if err != tc.err && err.Error() != tc.err.Error() {
			t.Errorf("retrieveElement(%v, %v): expected %v, got: %v", tc.json, tc.element, tc.err, err)
		}
	}
}
