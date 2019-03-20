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
		log.Fatalf("error: yaml: %v", err)
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
