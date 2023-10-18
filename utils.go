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
	"os"
)

// AppendHostsFile appends a string to /etc/hosts as a new line
func AppendHostsFile(content string) error {
	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fmt.Printf(fmt.Sprintf("\n%s\n", content))

	if _, err := f.Write([]byte(fmt.Sprintf("\n%s\n", content))); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
