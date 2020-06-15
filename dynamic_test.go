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
		actual, err := formatKey(tc.input, tc.encrypted)
		if actual != tc.expected {
			t.Errorf("formatKey(%v, %v): expected %v, actual %v", tc.input, tc.encrypted, tc.expected, actual)
		}
		if err != tc.err {
			t.Errorf("formatKey(%v, %v): expected %v, got: %v", tc.input, tc.encrypted, tc.err, err)
		}
	}
}

func TestValidateSignStringRS256PKCS8(t *testing.T) {
	var tests = []struct {
		args  []string
		valid bool
	}{
		{[]string{}, false},
		{[]string{"key"}, false},
		{[]string{"key", "passphrase"}, false},
		{[]string{"key", ""}, false},
		{[]string{"key", "passphrase", "string"}, true},
		{[]string{"key", "", "string"}, true},
	}

	for _, tc := range tests {
		actual := validateSignStringRS256PKCS8(tc.args)
		if actual != tc.valid {
			t.Errorf("validateSignStringRS256PKCS8(%s): expected %v, actual %v", tc.args, tc.valid, actual)
		}
	}
}

func TestArgsToStringToSign(t *testing.T) {
	var tests = []struct {
		existingHeaders map[string]string
		args            []string
		expected        string
	}{
		{
			map[string]string{},
			[]string{""},
			"\n",
		},
		{
			map[string]string{},
			[]string{"my string to sign"},
			"my string to sign\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"my string to sign"},
			"my string to sign\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header"},
			"123\n",
		},
		{
			map[string]string{
				"x-previous-header": "123",
			},
			[]string{"x-previous-header", "my string"},
			"123\nmy string\n",
		},
	}

	for _, tc := range tests {
		actual := argsToStringToSign(tc.existingHeaders, tc.args)
		if actual != tc.expected {
			t.Errorf("argsToStringToSign(%s, %s): expected %s, actual %s", tc.existingHeaders, tc.args, tc.expected, actual)
		}
	}
}
