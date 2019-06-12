// MIT License
//
// Copyright (c) 2018 Yunzhu Li
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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPRequestConfig type
type HTTPRequestConfig struct {
	Method               string
	URL                  string
	QueryParams          map[string]string
	Headers              map[string]string
	BasicAuthUsername    string
	BasicAuthPassword    string
	Body                 io.Reader
	Attempts             int
	TimeoutSeconds       time.Duration
	SkipCertVerification bool
}

// SendHTTPRequest sends an HTTP request and returns response body and status
func SendHTTPRequest(config *HTTPRequestConfig) (*http.Response, []byte, error) {
	// Check input
	if config == nil {
		return nil, nil, fmt.Errorf("config is nil")
	}

	if len(config.Method) <= 0 {
		return nil, nil, fmt.Errorf("Method is required")
	}

	if len(config.URL) <= 0 {
		return nil, nil, fmt.Errorf("URL is required")
	}

	if config.Attempts == 0 {
		config.Attempts = 1
	}

	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 5
	}

	// Create request
	req, err := http.NewRequest(
		config.Method,
		config.URL,
		config.Body,
	)

	if err != nil {
		return nil, nil, err
	}

	// Query params
	if len(config.QueryParams) > 0 {
		q := req.URL.Query()
		for k, v := range config.QueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// BasicAuth header
	if len(config.BasicAuthUsername) > 0 || len(config.BasicAuthPassword) > 0 {
		req.SetBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword)
	}

	// Other headers
	for k, v := range config.Headers {
		req.Header.Add(k, v)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SkipCertVerification},
	}

	client := http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.TimeoutSeconds * time.Second),
	}

	// Start sending request
	var resp *http.Response
	for a := config.Attempts; a > 0; a-- {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, nil, err
	}

	// Release resource when done
	defer resp.Body.Close()

	// Read body into a buffer
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, buf.Bytes(), nil
}
