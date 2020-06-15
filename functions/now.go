package functions

import (
	"strconv"
	"time"
)

// Now returns the number of seconds since the Unix epoch.
func Now(existingHeaders map[string]string, args []string) (string, error) {
	return strconv.FormatInt(time.Now().Unix(), 10), nil
}
