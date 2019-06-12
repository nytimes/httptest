// MIT License
//
// Copyright (c) 2019 Yunzhu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// TestResult stores results of a single test
type TestResult struct {
	Skipped bool
	Errors  []error
}

// GenerateTestInfoString generates a human-readable string indicates the test
func GenerateTestInfoString(test *Test) string {
	return fmt.Sprintf("%s | %s | [%s]", test.Filename, test.Description, test.Request.Path)
}

// RunTest runs single test
func RunTest(test *Test, defaultHost string) *TestResult {
	result := &TestResult{}

	if err := preProcessTest(test, defaultHost); err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}

	conditionsMet, err := validateConditions(test)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}
	if !conditionsMet {
		// Skip test
		result.Skipped = true
		return result
	}

	url := test.Request.Scheme + "://" + test.Request.Host + path.Join("/", test.Request.Path)

	var body io.Reader
	if len(test.Request.Body) > 0 {
		body = strings.NewReader(test.Request.Body)
	}

	reqConfig := &HTTPRequestConfig{
		Method:               test.Request.Method,
		URL:                  url,
		Headers:              test.Request.Headers,
		Body:                 body,
		Attempts:             1,
		TimeoutSeconds:       5,
		SkipCertVerification: test.SkipCertVerification,
	}

	resp, respBody, err := SendHTTPRequest(reqConfig)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}

	// Append response validation errors
	result.Errors = append(result.Errors, validateResponse(test, resp, respBody)...)

	return result
}

func preProcessTest(test *Test, defaultHost string) error {
	// Scheme
	scheme := stringValue(test.Request.Scheme, "https")
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("invalid scheme %s. only http and https are supported", scheme)
	}
	test.Request.Scheme = scheme

	// Host
	host := stringValue(test.Request.Host, defaultHost)
	if len(host) == 0 {
		return fmt.Errorf("no host specified for this test and no default host set")
	}
	test.Request.Host = host

	// Method
	method := stringValue(test.Request.Method, "GET")
	if method != "GET" && method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
		return fmt.Errorf("invalid method %s. only GET, POST, PUT, PATCH, DELETE are supported", method)
	}
	test.Request.Method = method

	// Path
	if len(test.Request.Path) < 1 {
		return fmt.Errorf("request path is required")
	}

	return nil
}

func stringValue(val, defaultVal string) string {
	if len(val) > 0 {
		return val
	}
	return defaultVal
}

func validateConditions(test *Test) (bool, error) {
	// Environment variable
	for key, pattern := range test.Conditions.Env {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return false, fmt.Errorf("%s", err.Error())
		}

		if !re.MatchString(os.Getenv(key)) {
			return false, nil
		}
	}

	return true, nil
}

func validateResponse(test *Test, response *http.Response, body []byte) []error {
	errors := []error{}

	errors = append(errors, validateResponseStatus(test, response)...)
	errors = append(errors, validateResponseHeaders(test, response)...)
	errors = append(errors, validateResponseBody(test, response, body)...)

	return errors
}

func validateResponseStatus(test *Test, response *http.Response) []error {
	errors := []error{}
	expected := test.Response

	matched := false
	for _, code := range expected.StatusCodes {
		if code == response.StatusCode {
			matched = true
		}
	}

	if !matched && len(expected.StatusCodes) > 0 {
		errors = append(errors, fmt.Errorf("unexpected status code - expected %v, got %d", expected.StatusCodes, response.StatusCode))
	}

	return errors
}

func validateResponseHeaders(test *Test, response *http.Response) []error {
	errors := []error{}
	expectedResponse := test.Response

	// Patterns
	patterns := expectedResponse.Headers.Patterns
	for header, pattern := range patterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s", err.Error()))
			continue
		}

		value := strings.ToLower(response.Header.Get(header))
		if !re.MatchString(value) {
			errors = append(errors, fmt.Errorf("response header \"%s\" has value \"%s\", does not match pattern \"%s\"", header, value, pattern))
		}
	}

	// NotPresent assertions
	npAssertions := expectedResponse.Headers.NotPresent
	for _, header := range npAssertions {
		if len(response.Header.Get(header)) > 0 {
			errors = append(errors, fmt.Errorf("found unexpected response header \"%s\"", header))
		}
	}

	return errors
}

func validateResponseBody(test *Test, response *http.Response, body []byte) []error {
	errors := []error{}

	patterns := test.Response.Body.Patterns
	for _, pattern := range patterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s", err.Error()))
			continue
		}

		if !re.Match(body) {
			errors = append(errors, fmt.Errorf("response body does not match pattern \"%s\"", pattern))
		}
	}

	return errors
}
