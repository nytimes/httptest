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
	"sync"
)

// RunTests runs all tests
func RunTests(tests []*Test, config *Config) bool {
	sem := make(chan byte, config.Concurrency)
	mux := sync.Mutex{}

	passed, failed, skipped := 0, 0, 0

	for _, test := range tests {
		// Attempt to take a slot
		sem <- 0

		go func(t *Test) {
			// Release a slot when done
			defer func() { <-sem }()

			result := RunTest(t, config.Host)

			// Acquire lock before accessing shared variables and writing output.
			// Code in critical section should not perform network I/O.
			mux.Lock()

			if result.Skipped {
				skipped++
				if !config.PrintFailedTestsOnly {
					PrintTestResult(t, result)
				}
			} else if len(result.Errors) > 0 {
				failed++
				PrintTestResult(t, result)
			} else {
				passed++
				if !config.PrintFailedTestsOnly {
					PrintTestResult(t, result)
				}
			}
			mux.Unlock()
		}(test)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- 0
	}

	PrintTestSummary(passed, failed, skipped)

	if failed > 0 {
		return false
	}
	return true
}
