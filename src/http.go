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

package src

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
		config.TimeoutSeconds = 10
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

	// Host header
	// https://github.com/golang/go/issues/7682
	if len(config.Headers["host"]) > 0 {
		req.Host = config.Headers["host"]
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
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
