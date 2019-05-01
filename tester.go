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
	"net/http"
	"path"
)

// RunTest runs single test
func RunTest(test *Test, defaultAddress string) error {
	scheme := stringValue(test.Request.Scheme, "https")
	address := stringValue(test.Request.Address, defaultAddress)
	if len(address) == 0 {
		return fmt.Errorf("no address specified for test %s, %s", test.Request.Path, test.Description)
	}
	method := stringValue(test.Request.Method, "GET")

	url := scheme + "://" + address + path.Join("/", test.Request.Path)

	reqConfig := &HTTPRequestConfig{
		Method:         method,
		URL:            url,
		Headers:        test.Request.Headers,
		Attempts:       1,
		TimeoutSeconds: 5,
	}

	resp, _, err := SendHTTPRequest(reqConfig)
	if err != nil {
		return err
	}

	return validateResponse(test, resp)
}

func validateResponse(test *Test, response *http.Response) error {
	tr := test.Response
	if tr.Status != 0 && tr.Status != response.StatusCode {
		return fmt.Errorf("unexpected status code. expecting %d, got %d", response.StatusCode, tr.Status)
	}

	return nil
}

func stringValue(val, defaultVal string) string {
	if len(val) > 0 {
		return val
	}
	return defaultVal
}
