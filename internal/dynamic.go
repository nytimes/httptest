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

package internal

import (
	"fmt"

	"github.com/nytimes/httptest/functions"
)

// ProcessDynamicHeaders creates headers based on the function and adds them to the map of all headers.
func ProcessDynamicHeaders(dynamicHeaders []DynamicHeader, allHeaders map[string]string) error {
	for _, dynamicHeader := range dynamicHeaders {
		if _, exists := allHeaders[dynamicHeader.Name]; exists {
			return fmt.Errorf("cannot process dynamic header %s; a header with that name is already defined", dynamicHeader.Name)
		}

		dhfn, exists := funcMap[dynamicHeader.Function]

		if !exists {
			return fmt.Errorf("unknown function %s", dynamicHeader.Function)
		}

		value, err := dhfn(allHeaders, dynamicHeader.Args)

		if err != nil {
			return err
		}

		allHeaders[dynamicHeader.Name] = value
	}

	return nil
}

// Generic signature for any function that can resolve a dynamic header value.
type resolveHeader func(existingHeaders map[string]string, args []string) (string, error)

// Map of strings to dynamic header functions.
var funcMap = map[string]resolveHeader{
	"now":                  functions.Now,
	"signStringRS256PKCS8": functions.SignStringRS256PKCS8,
	"postFormURLEncoded":   functions.PostFormURLEncoded,
	"concat":               functions.Concat,
}
