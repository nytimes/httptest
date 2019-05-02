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
	"sync"
)

// RunTests runs all tests
func RunTests(tests []*Test, config *Config) bool {
	mux := sync.Mutex{}
	sem := make(chan byte, config.Concurrency)

	passed, failed := 0, 0

	for _, test := range tests {
		// Attempt to take a slot
		sem <- 0

		go func(t *Test) {
			// Release a slot
			defer func() { <-sem }()

			err := RunTest(t, config.Host)

			// Acquire lock before accessing shared variables and writing output
			// Code in critical section should not perform network I/O
			mux.Lock()
			testInfoString := GenerateTestInfoString(t)
			if err != nil {
				failed++
				fmt.Printf("failed: %s\n%s\n\n", testInfoString, err)
			} else {
				if config.PrintSuccessfulTests {
					fmt.Printf("passed: %s\n", testInfoString)
				}
				passed++
			}
			mux.Unlock()
		}(test)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- 0
	}

	fmt.Printf("\n%d/%d tests passed\n", passed, passed+failed)
	if failed > 0 {
		return false
	}
	return true
}
