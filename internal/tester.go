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

package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"

	methodGET     = "GET"
	methodPOST    = "POST"
	methodPUT     = "PUT"
	methodPATCH   = "PATCH"
	methodDELETE  = "DELETE"
	methodHEAD    = "HEAD"
	methodOPTIONS = "OPTIONS"
	methodPURGE   = "PURGE"
	methodPROPFIND = "PROPFIND"
)

type TestResult struct {
	Retries int
	Skipped bool
	Errors  []error
}

func RunTest(test *Test, defaultHost string, maxRetries int) *TestResult {
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

	url := test.Request.Scheme + "://" + test.Request.Host + test.Request.Path

	var body io.Reader
	if len(test.Request.Body) > 0 {
		body = strings.NewReader(test.Request.Body)
	}

	retryCallback := func(ctx context.Context, resp *http.Response, inErr error) (bool, error) {
		if inErr != nil {
			// retry if there is an error with the request
			return true, nil
		}

		errs := validateResponseStatus(test, resp)
		if len(errs) >= 1 {
			// retry if there is an error
			return true, nil
		}

		// stop retrying
		return false, nil
	}

	reqConfig := &HTTPRequestConfig{
		Method:               test.Request.Method,
		URL:                  url,
		Headers:              test.Request.Headers,
		Body:                 body,
		TimeoutSeconds:       60,
		SkipCertVerification: test.SkipCertVerification,
		RetryCallback:        retryCallback,
		MaxRetries:           maxRetries,
	}

	zap.L().Info("sending request",
		zap.Any("request", reqConfig),
	)

	for i := 0; i <= maxRetries; i++ {
		result.Errors = []error{}
		result.Retries = i

		resp, respBody, err := SendHTTPRequest(reqConfig)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}

		zap.L().Info("got response",
			zap.ByteString("body", respBody),
			zap.String("status", resp.Status),
			zap.Any("headers", resp.Header),
		)

		// Append response validation errors
		result.Errors = append(result.Errors, validateResponse(test, resp, respBody)...)

		if len(result.Errors) == 0 {
			return result
		}
	}

	return result
}

func preProcessTest(test *Test, defaultHost string) error {
	// Scheme
	scheme := stringValue(test.Request.Scheme, schemeHTTPS)
	switch scheme {
	case schemeHTTP, schemeHTTPS:
		test.Request.Scheme = scheme
	default:
		return fmt.Errorf("invalid scheme %s. only http and https are supported", scheme)
	}

	host := stringValue(test.Request.Host, defaultHost)
	if len(host) == 0 {
		return fmt.Errorf("no host specified for this test and no default host set")
	}
	test.Request.Host = host

	method := stringValue(test.Request.Method, methodGET)
	switch method {
	case methodGET, methodPOST, methodPUT, methodPATCH, methodDELETE, methodHEAD, methodOPTIONS, methodPURGE, methodPROPFIND:
		test.Request.Method = method
	default:
		return fmt.Errorf("invalid method %s. only GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, PURGE, PROPFIND are supported", method)
	}

	if len(test.Request.Path) < 1 {
		return fmt.Errorf("request path is required")
	}

	if !strings.HasPrefix(test.Request.Path, "/") {
		return fmt.Errorf("request.path must start with /")
	}

	if err := ProcessDynamicHeaders(test.Request.DynamicHeaders, test.Request.Headers); err != nil {
		return err
	}

	// Convert header fields to lowercase
	// https://tools.ietf.org/html/rfc7540#section-8.1.2
	headers := make(map[string]string, len(test.Request.Headers))
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
		re, err := compilePattern(pattern)
		if err != nil {
			return false, err
		}

		if !re.MatchString(os.Getenv(key)) {
			return false, nil
		}
	}

	return true, nil
}

func validateResponse(test *Test, response *http.Response, body []byte) (errs []error) {
	errs = append(errs, validateResponseStatus(test, response)...)
	errs = append(errs, validateResponseHeaders(test, response)...)
	errs = append(errs, validateResponseBody(test, response, body)...)
	return errs
}

func validateResponseStatus(test *Test, response *http.Response) []error {
	errors := []error{}
	expected := test.Response

	matched := false
	for _, code := range expected.StatusCodes {
		if code == response.StatusCode {
			matched = true
			break
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

	errors = append(errors, validateHeaderPatterns(response, expectedResponse.Headers.Patterns, true)...)

	// NotMatching assertions
	errors = append(errors, validateHeaderPatterns(response, expectedResponse.Headers.NotMatching, false)...)

	errors = append(errors, validateHeadersNotPresent(response, expectedResponse.Headers.NotPresent)...)

	errors = append(errors, validateHeadersIfPresentNotMatching(response, expectedResponse.Headers.IfPresentNotMatching)...)

	return errors
}

func validateHeaderPatterns(response *http.Response, patterns map[string]string, expectedToMatch bool) []error {
	errors := []error{}

	for header, pattern := range patterns {
		re, err := compilePattern(pattern)
		if err != nil {
			errors = append(errors, fmt.Errorf("invalid test pattern `%s`: %s", pattern, err.Error()))
			continue
		}

		value := response.Header.Get(header)
		if value == "" {
			if expectedToMatch {
				errors = append(errors, fmt.Errorf("response header \"%s\" not found, expected to match pattern \"%s\"", header, pattern))
			} else {
				errors = append(errors, fmt.Errorf("response header \"%s\" not found, expected to be present", header))
			}
			continue
		}

		if expectedToMatch && !re.MatchString(value) {
			errors = append(errors, fmt.Errorf("response header \"%s\" has value \"%s\", which does not match pattern \"%s\"", header, value, pattern))
		}

		if !expectedToMatch && re.MatchString(value) {
			errors = append(errors, fmt.Errorf("response header \"%s\" has value \"%s\", which matches pattern \"%s\"", header, value, pattern))
		}
	}

	return errors
}

func validateHeadersNotPresent(response *http.Response, notPresentAssertions []string) []error {
	errors := []error{}

	for _, header := range notPresentAssertions {
		if len(response.Header.Get(header)) > 0 {
			errors = append(errors, fmt.Errorf("found unexpected response header \"%s\"", header))
		}
	}

	return errors
}

func validateHeadersIfPresentNotMatching(response *http.Response, ifPresentNotMatchingAssertions map[string]string) []error {
	errors := []error{}

	for header := range ifPresentNotMatchingAssertions {
		if response.Header.Get(header) != "" {
			errors = append(errors, validateHeaderPatterns(response, ifPresentNotMatchingAssertions, false)...)
		}
	}

	return errors
}

func validateResponseBody(test *Test, response *http.Response, body []byte) []error {
	errors := []error{}

	patterns := test.Response.Body.Patterns
	for _, pattern := range patterns {
		re, err := compilePattern(pattern)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if !re.Match(body) {
			errors = append(errors, fmt.Errorf("response body does not match pattern \"%s\"", pattern))
		}
	}

	return errors
}

func compilePattern(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile("(?i)" + pattern)
}
