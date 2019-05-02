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
	"log"
	"os"
	"strconv"
)

var (
	// BuildBranch is the branch of the source the binary built from
	BuildBranch string

	// BuildCommit is the commit hash
	BuildCommit string

	// BuildTime is the build time
	BuildTime string
)

// Config stores application configuration
type Config struct {
	Concurrency          int
	Host                 string
	DNSOverride          string
	PrintSuccessfulTests bool
}

// FromEnv returns config read from environment variables
func FromEnv() *Config {
	// Parse non-string values
	concurrency, err := strconv.Atoi(getEnv("TEST_CONCURRENCY", "2"))
	if err != nil {
		log.Fatalf("invalid concurrency value: %s", err)
	}

	if concurrency < 1 {
		log.Fatalf("invalid concurrency value: %d", concurrency)
	}

	printSuccessfulTests := true
	if getEnv("TEST_PRINT_SUCCESSFUL_TESTS", "true") == "false" {
		printSuccessfulTests = false
	}

	return &Config{
		Concurrency:          concurrency,
		Host:                 getEnv("TEST_HOST", ""),
		DNSOverride:          getEnv("TEST_DNS_OVERRIDE", ""),
		PrintSuccessfulTests: printSuccessfulTests,
	}
}

// Read environment variable with default values
func getEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}
