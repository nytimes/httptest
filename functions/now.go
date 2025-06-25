package functions

import (
	"strconv"
	"time"
)

// Now returns the number of seconds since the Unix epoch.
func Now(_ map[string]string, _ []string) (string, error) {
	return strconv.FormatInt(time.Now().Unix(), 10), nil
}

// was func Now(existingHeaders map[string]string, args []string) (string, error) {
// but we apparently never did anything with existingHeaders or args
