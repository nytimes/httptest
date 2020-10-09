package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const timeoutSeconds = 10 // the max number of seconds to wait for a response
var errURLMissing = errors.New("error calling PostData; URL is required")
var errURLEmpty = errors.New("error calling PostData; URL is empty")

// PostData sends an HTTP POST request to a given URL with the data provided,
// returning either the whole response body or a specified JSON element.
func PostData(existingHeaders map[string]string, args []string) (string, error) {
	// Ensure that we have at least one argument for the URL
	if len(args) == 0 {
		return "", errURLMissing
	}

	// Get the URL, response element, and request body from the args
	url, responseElement, err := argsToURLResponseElement(args)
	if err != nil {
		return "", err
	}
	var requestBody string
	if len(args) > 2 {
		requestBody = argsToRequestBody(existingHeaders, args[2:])
	}

	// Send the request
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	client.Timeout = time.Second * timeoutSeconds
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Get the specific element requested from the response
	var result string
	jsonBody := make(map[string]json.RawMessage)
	err = json.Unmarshal(responseBody, &jsonBody)
	if err != nil {
		// This must be an array, so just return the whole response since we can't determine which element to return
		responseElement = ""
	}
	if responseElement != "" {
		// Return the value associated with the key responseElement,
		// either a plain string for primitives, or a JSON string for objects/arrays
		value := jsonBody[responseElement]
		if value == nil {
			return "", fmt.Errorf("error: element %s not found in response body", responseElement)
		}
		str, err := jsonToString(value)
		if err != nil {
			return "", err
		}
		result = str
	} else {
		// No specific element was requested, so return the whole response body as a JSON string
		str, err := jsonToString(responseBody)
		if err != nil {
			return "", err
		}
		result = str
	}

	return result, nil
}

// Gets the URL and (optional) responseElement from the args.
func argsToURLResponseElement(args []string) (string, string, error) {
	if args[0] == "" {
		return "", "", errURLEmpty
	}

	url := args[0]
	responseElement := ""

	if len(args) > 1 {
		responseElement = args[1]
	}

	return url, responseElement, nil
}

// Concatenates the rest of the args (literal strings or existing header values), delimiting with '&'.
func argsToRequestBody(existingHeaders map[string]string, args []string) string {
	var buffer strings.Builder

	for index, arg := range args {
		if value, ok := existingHeaders[arg]; ok {
			buffer.WriteString(value)
		} else {
			buffer.WriteString(arg)
		}
		if index < (len(args) - 1) {
			buffer.WriteString("&")
		}
	}

	return buffer.String()
}

// Returns whether the string represented by the byte slice is a JSON object
func isJSONObject(str []byte) bool {
	test := make(map[string]json.RawMessage)
	err := json.Unmarshal(str, &test)
	return err == nil
}

// Returns whether the string represented by the byte slice is a JSON array
func isJSONArray(str []byte) bool {
	var test []interface{}
	err := json.Unmarshal(str, &test)
	return err == nil
}

// Takes raw encoded JSON and converts to string
func jsonToString(rawJSON []byte) (string, error) {
	if isJSONObject(rawJSON) || isJSONArray(rawJSON) {
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, rawJSON); err != nil {
			return "", err
		}
		return buffer.String(), nil
	} else {
		var primitiveValue interface{}
		if err := json.Unmarshal(rawJSON, &primitiveValue); err != nil {
			return "", err
		}
		str := fmt.Sprintf("%v", primitiveValue)
		return str, nil
	}
}
