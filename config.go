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

package httptest

import (
	"fmt"
	"os"
	"strconv"
)

// Config stores application configuration
type Config struct {
	Concurrency          int
	DNSOverride          string
	Host                 string
	PrintFailedTestsOnly bool
	TestDirectory        string
}

// FromEnv returns config read from environment variables
func FromEnv() (*Config, error) {
	// Parse non-string values
	concurrency, err := strconv.Atoi(getEnv("TEST_CONCURRENCY", "2"))
	if err != nil {
		return nil, fmt.Errorf("invalid concurrency value: %s", err)
	}
	if concurrency < 1 {
		return nil, fmt.Errorf("invalid concurrency value: %d", concurrency)
	}

	printFailedOnly := false
	if getEnv("TEST_PRINT_FAILED_ONLY", "false") == "true" {
		printFailedOnly = true
	}

	return &Config{
		Concurrency:          concurrency,
		Host:                 getEnv("TEST_HOST", ""),
		DNSOverride:          getEnv("TEST_DNS_OVERRIDE", ""),
		PrintFailedTestsOnly: printFailedOnly,
		TestDirectory:        getEnv("TEST_DIRECTORY", "tests"),
	}, nil
}

// ApplyConfig applies config
func ApplyConfig(config *Config) error {
	if len(config.DNSOverride) > 0 {
		if len(config.Host) < 1 {
			return fmt.Errorf("TEST_HOST is required to use DNS override")
		}
		return AppendHostsFile(fmt.Sprintf("%s %s", config.DNSOverride, config.Host))
	}
	return nil
}

// Read environment variable with default values
func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}
