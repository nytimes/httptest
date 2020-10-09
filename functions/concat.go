package functions

import (
	"strings"
)

// Concat concatenates the args into a single string, substituting previously defined headers if available.
func Concat(existingHeaders map[string]string, args []string) (string, error) {
	var buffer strings.Builder

	for _, arg := range args {
		if value, ok := existingHeaders[arg]; ok {
			buffer.WriteString(value)
		} else {
			buffer.WriteString(arg)
		}
	}

	return buffer.String(), nil
}
