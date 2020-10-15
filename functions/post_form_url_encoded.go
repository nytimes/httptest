package functions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
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

	return retrieveElement(responseBody, responseElement)
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

func retrieveElement(data []byte, field string) (string, error) {
	if !gjson.ValidBytes(data) {
		return "", fmt.Errorf("invalid JSON")
	}

	var result string
	if field == "" {
		result = string(data)
	} else {
		result = gjson.GetBytes(data, field).String()
	}

	return string(pretty.Ugly([]byte(result))), nil
}
