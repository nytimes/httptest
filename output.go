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

	"github.com/fatih/color"
)

// PrintTestResult prints result of a single test
func PrintTestResult(test *Test, result *TestResult) {
	fmt.Println("")

	// Print colored status text
	if result.Skipped {
		color.HiBlue("SKIPPED")
	} else if len(result.Errors) < 1 {
		color.HiGreen("PASSED")
	} else {
		color.HiRed("FAILED")
	}

	// Print test info
	fmt.Printf("%s\n", fmt.Sprintf("%s | %s | [%s]", test.Filename, test.Description, test.Request.Path))

	// Print all errors
	if len(result.Errors) > 0 {
		fmt.Printf("errors:\n")
		for _, err := range result.Errors {
			fmt.Printf("%s\n", err.Error())
		}
	}
}

// PrintTestSummary prints summary info for multiple tests
func PrintTestSummary(passed, failed, skipped int) {
	fmt.Printf(
		"\n%s passed\n%s failed\n%s skipped\n",
		color.HiGreenString("%d", passed),
		color.HiRedString("%d", failed),
		color.HiBlueString("%d", skipped),
	)
}
