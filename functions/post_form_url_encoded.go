package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const timeoutSeconds = 10 // the max number of seconds to wait for a response
var errURLMissing = errors.New("error calling PostFormURLEncoded; URL is required")
var errURLEmpty = errors.New("error calling PostFormURLEncoded; URL is empty")

// PostFormURLEncoded sends an HTTP POST request to a given URL with the data provided,
// returning either the whole response body or a specified JSON element.
func PostFormURLEncoded(existingHeaders map[string]string, args []string) (string, error) {
	// Get the URL, response element, and request body from the args
	if len(args) == 0 {
		return "", errURLMissing
	}

	if args[0] == "" {
		return "", errURLEmpty
	}

	endpoint := args[0]

	var responseElement string
	if len(args) > 1 {
		responseElement = args[1]
	}

	var requestBody url.Values
	if len(args) > 2 {
		requestBody = argsToRequestBody(existingHeaders, args[2:])
	}

	// Send the request
	client := &http.Client{}
	client.Timeout = time.Second * timeoutSeconds
	response, err := client.PostForm(endpoint, requestBody)
	if err != nil {
		return "", err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Get the specific element requested from the response
	result, err := retrieveElementTree(responseBody, responseElement)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// Takes the rest of the args (literal strings or existing header values) and transforms them into a map of form values.
func argsToRequestBody(existingHeaders map[string]string, args []string) url.Values {
	values := url.Values{}

	for _, arg := range args {
		var parts []string
		if value, ok := existingHeaders[arg]; ok {
			parts = strings.SplitN(value, "=", 2)
		} else {
			parts = strings.SplitN(arg, "=", 2)
		}
		values.Add(parts[0], parts[1])
	}

	return values
}

// Retrieves the JSON element (specified by a dot-delimited string)
func retrieveElementTree(jsonContents []byte, str string) ([]byte, error) {
	elements := strings.Split(str, ".")
	var err error
	for _, element := range elements {
		jsonContents, err = retrieveElement(jsonContents, element)
		if err != nil {
			return nil, err
		}
	}

	return jsonContents, nil
}

// Retrieve the element by parsing into a map for a JSON object or slice for a JSON array.
func retrieveElement(jsonContents []byte, element string) ([]byte, error) {
	// Parse into a map for a JSON object
	jsonObject := make(map[string]json.RawMessage)
	err := json.Unmarshal(jsonContents, &jsonObject)
	if err != nil {
		// It must be a JSON array, so parse into a slice
		jsonArray := make([]json.RawMessage, 0)
		err = json.Unmarshal(jsonContents, &jsonArray)
		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(element)
		if err != nil {
			return nil, fmt.Errorf("invalid index for JSON array: %v", element)
		}

		if len(jsonArray) <= i || i < 0 {
			return nil, fmt.Errorf("JSON array index out of bounds: %v", element)
		}

		return jsonToBytes(jsonArray[i])
	}

	if value, ok := jsonObject[element]; ok {
		return jsonToBytes(value)
	}

	return nil, fmt.Errorf("error: element %s not found in response body", element)
}

// Transforms raw encoded JSON into a canonical format (no whitespace) and converts it to a byte slice
func jsonToBytes(rawJSON []byte) ([]byte, error) {
	if isJSONObject(rawJSON) || isJSONArray(rawJSON) {
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, rawJSON); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	var primitiveValue interface{}
	if err := json.Unmarshal(rawJSON, &primitiveValue); err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%v", primitiveValue)
	return []byte(str), nil
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
