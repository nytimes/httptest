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

// GenerateTestInfoString generates a human-readable string indicates the test
func GenerateTestInfoString(test *Test) string {
	return fmt.Sprintf("%s | [%s]", test.Description, test.Request.Path)
}

// RunTest runs single test
func RunTest(test *Test, defaultHost string) error {
	if err := preProcessTest(test, defaultHost); err != nil {
		return err
	}

	conditionsMet, err := validateConditions(test)
	if err != nil {
		return err
	}
	if !conditionsMet {
		// Skip test
		return nil
	}

	url := test.Request.Scheme + "://" + test.Request.Host + path.Join("/", test.Request.Path)

	var body io.Reader
	if len(test.Request.Body) > 0 {
		body = strings.NewReader(test.Request.Body)
	}

	reqConfig := &HTTPRequestConfig{
		Method:         test.Request.Method,
		URL:            url,
		Headers:        test.Request.Headers,
		Body:           body,
		Attempts:       1,
		TimeoutSeconds: 5,
	}

	resp, respBody, err := SendHTTPRequest(reqConfig)
	if err != nil {
		return err
	}

	return validateResponse(test, resp, respBody)
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
	if method != "GET" && method != "POST" {
		return fmt.Errorf("invalid method %s. only GET and POST are supported", method)
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

func validateResponse(test *Test, response *http.Response, body []byte) error {
	if err := validateResponseStatus(test, response); err != nil {
		return err
	}

	if err := validateResponseHeaders(test, response); err != nil {
		return err
	}

	if err := validateResponseBody(test, response, body); err != nil {
		return err
	}

	return nil
}

func validateResponseStatus(test *Test, response *http.Response) error {
	expected := test.Response
	for _, code := range expected.StatusCodes {
		if code != response.StatusCode {
			return fmt.Errorf("unexpected status code. expected %v, got %d", expected.StatusCodes, response.StatusCode)
		}
	}
	return nil
}

func validateResponseHeaders(test *Test, response *http.Response) error {
	expectedResponse := test.Response

	// Patterns
	patterns := expectedResponse.Headers.Patterns
	for header, pattern := range patterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}

		value := strings.ToLower(response.Header.Get(header))
		if !re.MatchString(value) {
			return fmt.Errorf("the value of response header \"%s: %s\" does not match pattern \"%s\"", header, value, pattern)
		}
	}

	// NotPresent assertions
	npAssertions := expectedResponse.Headers.NotPresent
	for _, header := range npAssertions {
		if len(response.Header.Get(header)) > 0 {
			return fmt.Errorf("found unexpected response header \"%s\"", header)
		}
	}

	return nil
}

func validateResponseBody(test *Test, response *http.Response, body []byte) error {
	patterns := test.Response.Body.Patterns
	for _, pattern := range patterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}

		if !re.Match(body) {
			return fmt.Errorf("response body does not match pattern \"%s\"", pattern)
		}
	}
	return nil
}
