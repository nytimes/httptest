package httptest

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
	var formatKeyTests = []struct {
		input       string
		encrypted   bool
		expected    string
		expectError bool
	}{
		{unencrypted, false, unencryptedExpected, false},
		{encrypted, true, encryptedExpected, false},
		{unencrypted, true, "", true},
		{encrypted, false, "", true},
	}
	for _, tc := range formatKeyTests {
		actual, err := FormatKey(tc.input, tc.encrypted)
		if actual != tc.expected {
			t.Errorf("FormatKey(%v, %v): expected %v, actual %v", tc.input, tc.encrypted, tc.expected, actual)
		}
		if tc.expectError && err == nil {
			t.Errorf("FormatKey(%v, %v): expected error, actual %v", tc.input, tc.encrypted, nil)
		}
		if !tc.expectError && err != nil {
			t.Errorf("FormatKey(%v, %v): expected no error, actual %v", tc.input, tc.encrypted, err)
		}
	}
}
