package main

import (
	"testing"
)

func TestFormatKey(t *testing.T) {
	unencrypted := `-----BEGIN PRIVATE KEY----- abcdef ghijkl mnopqr -----END PRIVATE KEY-----`
	unencryptedExpected := `-----BEGIN PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END PRIVATE KEY-----`

	encrypted := `-----BEGIN ENCRYPTED PRIVATE KEY----- abcdef ghijkl mnopqr -----END ENCRYPTED PRIVATE KEY-----`
	encryptedExpected := `-----BEGIN ENCRYPTED PRIVATE KEY-----
abcdef
ghijkl
mnopqr
-----END ENCRYPTED PRIVATE KEY-----`

	var FormatKeyTests = []struct {
		input     string
		encrypted bool
		expected  string
		err       error
	}{
		{unencrypted, false, unencryptedExpected, nil},
		{encrypted, true, encryptedExpected, nil},
		{unencrypted, true, "", errInvalidKeyFormat},
		{encrypted, false, "", errInvalidKeyFormat},
	}

	for _, tc := range FormatKeyTests {
		actual, err := FormatKey(tc.input, tc.encrypted)
		if actual != tc.expected {
			t.Errorf("FormatKey(%v, %v): expected %v, actual %v", tc.input, tc.encrypted, tc.expected, actual)
		}
		if err != tc.err {
			t.Errorf("FormatKey(%v, %v): expected %v, got: %v", tc.input, tc.encrypted, tc.err, err)
		}
	}
}
