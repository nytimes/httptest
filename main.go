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
	fmt.Printf("httptest: %s %s %s\n", BuildCommit, BuildBranch, BuildTime)

	// Get and apply config
	config, err := FromEnv()
	if err != nil {
		log.Fatalf("error: failed to parse config: %s", err)
	}

	if err := ApplyConfig(config); err != nil {
		log.Fatalf("error: failed to apply config: %s", err)
	}

	// Parse and run tests
	tests, err := ParseAllTestsInDirectory(config.TestDirectory)
	if err != nil {
		log.Fatalf("error: failed to parse tests: %s", err)
	}

	if !RunTests(tests, config) {
		os.Exit(1)
	}

	os.Exit(0)
}
