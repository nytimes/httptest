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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// TestConfig is a single test config file
type TestConfig struct {
	Tests []*Test `yaml:"tests"`
}

// Test is a single test
type Test struct {
	Description string
	Request     struct {
		Scheme  string            `yaml:"scheme"`
		Address string            `yaml:"address"`
		Method  string            `yaml:"method"`
		Path    string            `yaml:"path"`
		Headers map[string]string `yaml:"headers"`
		Body    string            `yaml:"body"`
	} `yaml:"request"`
	Response struct {
		Status  int `yaml:"status"`
		Headers struct {
			Patterns map[string]string `yaml:"patterns"`
			Exclude  []string          `yaml:"exclude"`
		} `yaml:"headers"`
		Body struct {
			Patterns []string `yaml:"patterns"`
		}
	} `yaml:"response"`
}

// ParseAllTestConfigsInDirectory recursively parses all test config files in a given directory
func ParseAllTestConfigsInDirectory(root string) []*Test {
	files := []string{}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("error: %s", err)
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

	tests := []*Test{}
	for _, p := range files {
		tests = append(tests, parseTestConfig(p)...)
	}

	return tests
}

func parseTestConfig(path string) []*Test {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("ioutil: %v", err)
	}

	tc := &TestConfig{}
	err = yaml.Unmarshal(data, tc)
	if err != nil {
		log.Fatalf("unable to parse file %s: %v", path, err)
	}
	return tc.Tests
}
