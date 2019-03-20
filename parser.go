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

type testConfig struct {
	Tests []test `yaml:"tests"`
}

type test struct {
	Description string
	Request     struct {
		Scheme  string            `yaml:"scheme"`
		Method  string            `yaml:"method"`
		Path    string            `yaml:"path"`
		Headers map[string]string `yaml:"headers"`
	} `yaml:"request"`
	Response struct {
		Status  int `yaml:"status"`
		Headers struct {
			Match   []map[string]string `yaml:"match"`
			Pattern []map[string]string `yaml:"pattern"`
			Exclude []map[string]string `yaml:"exclude"`
		} `yaml:"headers"`
	} `yaml:"response"`
}

func parseTestConfig(path string) []test {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: ioutil: %v", err)
	}

	tc := &testConfig{}
	err = yaml.Unmarshal(data, tc)
	if err != nil {
		log.Fatalf("unable to parse file %s: %v", path, err)
	}
	return tc.Tests
}

func parseAllTestConfigsInDirectory(root string) []test {
	files := []string{}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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

	tests := []test{}
	for _, p := range files {
		tests = append(tests, parseTestConfig(p)...)
	}

	return tests
}
