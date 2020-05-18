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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/drone/envsubst"
	"gopkg.in/yaml.v2"
)

// TestFile is a single test definition file
type TestFile struct {
	Tests []*Test `yaml:"tests"`
}

// Test is a single test
type Test struct {
	Filename    string
	Description string
	Conditions  struct {
		Env map[string]string `yaml:"env"`
	} `yaml:"conditions"`
	SkipCertVerification bool `yaml:"skipCertVerification"`
	Request              struct {
		Scheme  string                   `yaml:"scheme"`
		Host    string                   `yaml:"host"`
		Method  string                   `yaml:"method"`
		Path    string                   `yaml:"path"`
		Headers map[string]string        `yaml:"headers"`
		DynamicHeaders map[string]string `yaml:"dynamicHeaders"`
		Body    string                   `yaml:"body"`
	} `yaml:"request"`
	Response struct {
		StatusCodes []int `yaml:"statusCodes"`
		Headers     struct {
			Patterns   map[string]string `yaml:"patterns"`
			NotPresent []string          `yaml:"notPresent"`
			NotMatching map[string]string `yaml:"notMatching"`
		} `yaml:"headers"`
		Body struct {
			Patterns []string `yaml:"patterns"`
		}
	} `yaml:"response"`
}

// ParseAllTestsInDirectory recursively parses all test definition files in a given directory
func ParseAllTestsInDirectory(root string) ([]*Test, error) {
	files := []string{}
	var dirWalkError error

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("error: %s\n", err)
			dirWalkError = err
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if dirWalkError != nil {
		return nil, dirWalkError
	}

	allTests := []*Test{}
	for _, p := range files {
		tests, err := parseTestFile(p)
		if err != nil {
			return nil, err
		}
		allTests = append(allTests, tests...)
	}

	return allTests, nil
}

func parseTestFile(filePath string) ([]*Test, error) {
	// Read file into buffer
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ioutil: %v", err)
	}

	// Environment variable substitution
	yamlString, err := envsubst.EvalEnv(string(data))
	if err != nil {
		return nil, fmt.Errorf("unable to parse file %s: %v", filePath, err)
	}

	// Parse as YAML
	tf := TestFile{}
	err = yaml.Unmarshal([]byte(yamlString), &tf)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file %s: %v", filePath, err)
	}

	// Add file path to tests
	fileName := path.Base(filePath)
	for _, test := range tf.Tests {
		test.Filename = fileName
	}

	return tf.Tests, nil
}
