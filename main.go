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
	"log"
	"os"
)

var (
	// BuildBranch is the branch of the source the binary built from
	BuildBranch string

	// BuildCommit is the commit hash
	BuildCommit string

	// BuildTime is the build time
	BuildTime string
)

func main() {
	// Print version info
	fmt.Printf("httptest: %s %s %s\n\n", BuildCommit, BuildBranch, BuildTime)

	// Get and apply config
	config, err := FromEnv()
	if err != nil {
		log.Fatalf("error: failed to parse config: %s", err)
	}

	if err := ApplyConfig(config); err != nil {
		log.Fatalf("error: failed to apply config: %s", err)
	}

	// Parse and run tests
	tests, err := ParseAllTestsInDirectory("tests")
	if err != nil {
		log.Fatalf("error: failed to parse tests: %s", err)
	}

	if !RunTests(tests, config) {
		os.Exit(1)
	}

	os.Exit(0)
}
