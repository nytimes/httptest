// Copyright 2019 The New York Times Company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// TestResult stores results of a single test
type TestResult struct {
	Skipped bool
	Errors  []error
}

// RunTest runs single test
func RunTest(test *Test, defaultHost string) *TestResult {
	result := &TestResult{}

	// Validate test and assign default values
	if err := preProcessTest(test, defaultHost); err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}

	// Check test conditions and skip if not met
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

	url := test.Request.Scheme + "://" + test.Request.Host + test.Request.Path

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
		TimeoutSeconds:       60,
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

// preProcessTest validates test and assigns default values
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

	if !strings.HasPrefix(test.Request.Path, "/") {
		return fmt.Errorf("request.path must start with /")
	}

	// Convert header fields to lowercase
	// https://tools.ietf.org/html/rfc7540#section-8.1.2
	headers := map[string]string{}
	for k, v := range test.Request.Headers {
		headers[strings.ToLower(k)] = v
	}
	test.Request.Headers = headers

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

		// Get all instances of the response header
		values, ok := response.Header[http.CanonicalHeaderKey(header)]
		if !ok {
			errors = append(errors, fmt.Errorf("response header \"%s\" not found, expected to match pattern \"%s\"", header, pattern))
		}

		// Try to match pattern from one of the instances
		matched := false
		for _, value := range values {
			value = strings.ToLower(value)
			if re.MatchString(value) {
				matched = true
			}
		}

		if !matched {
			errors = append(errors, fmt.Errorf("response header \"%s\" has value(s) \"%s\", none of which match pattern \"%s\"", header, strings.Join(values[:], "\", \""), pattern))
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
